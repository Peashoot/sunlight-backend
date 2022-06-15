package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/peashoot/sunlight/config"
	"github.com/peashoot/sunlight/entity/db"
	"github.com/peashoot/sunlight/log"
	"github.com/peashoot/sunlight/utils"
	"gorm.io/gorm"
)

type LabelService struct {
}

func NewLabelService() *LabelService {
	return &LabelService{}
}

// AddOrGetExists 添加新的标签信息（标签的详情也具有唯一性）
func (labelService *LabelService) AddOrGetExists(labelInfo, userCode string) (db.UserLabelModel, error) {
	var label db.UserLabelModel
	locked := false
	defer func() {
		if locked {
			utils.RedisTryUnlock(utils.LabelAddKeyLockPrefix + labelInfo)
		}
	}()
	expiration := time.Now().Add(time.Duration(config.GetValue[int](config.RCN_LableWaitLockTimeOut)) * time.Millisecond).UnixMilli()
	for time.Now().UnixMilli() < expiration {
		locked, _ = utils.RedisTryLock(utils.LabelAddKeyLockPrefix+labelInfo,
			time.Duration(config.GetValue[int](config.RCN_LabelAddLockExpiration))*time.Second)
		if locked {
			break
		}
	}
	if !locked {
		return label, errors.New("failed to acquire lock")
	}
	// 如果缓存中存在，返回缓存中的实体
	label, err := utils.RedisGet[db.UserLabelModel](utils.UserLabelCachePrefix + labelInfo)
	if err == nil {
		return label, err
	}
	// 最后从数据库取
	err = config.MysqlDB.Where("name = ?", labelInfo).First(&label).Error
	if err == gorm.ErrRecordNotFound {
		// 添加记录到数据库中
		log.Info("[LabelService.AddOrGetExists]", "add new label", labelInfo, "to db")
		label = db.UserLabelModel{
			BaseModel: db.BaseModel{
				Code:      uuid.NewString(),
				CreatedBy: userCode,
			},
			Name:   labelInfo,
			Status: db.LabelExaminingStatus,
		}
		if err := config.MysqlDB.Create(&label).Error; err != nil {
			return label, err
		}
		// TODO: 将待审核的标签任务推送到审核员
		utils.RedisSetT(utils.UserLabelCachePrefix+labelInfo, label,
			time.Minute*time.Duration(config.GetValue[int](config.RCN_UserLabelCacheExpiration)))
		return label, nil
	}
	if err != nil {
		utils.RedisSetT(utils.UserLabelCachePrefix+labelInfo, label,
			time.Minute*time.Duration(config.GetValue[int](config.RCN_UserLabelCacheExpiration)))
	}
	return label, err
}

// Examine 审核标签
func (labelService *LabelService) Examine(labelCode, userCode string, status uint) (db.UserLabelModel, error) {
	// 修改标签审核状态
	log.Info("[LabelService.Examine]", "examine label", labelCode, "in db")
	var label db.UserLabelModel
	if err := config.MysqlDB.Where("code = ?", labelCode).First(&label).Error; err != nil {
		return label, err
	}
	label.Status = status
	label.UpdatedBy = userCode
	label.UpdatedAt = time.Now()
	if err := config.MysqlDB.Save(&label).Error; err != nil {
		return label, err
	}
	// 如果审核状态是不通过，则不添加到缓存中
	if status < 5 {
		return label, nil
	}
	// 添加到缓存
	utils.RedisSetT(utils.UserLabelCachePrefix+labelCode, label,
		time.Minute*time.Duration(config.GetValue[int](config.RCN_UserLabelCacheExpiration)))
	return label, nil
}

// Remove 移除标签
func (labelService *LabelService) Remove(labelCode, userCode string) (db.UserLabelModel, error) {
	// 从数据库删除
	log.Info("[LabelService.Remove]", "remove label", labelCode, "from db")
	var label db.UserLabelModel
	var err error
	if err = config.MysqlDB.Where("code = ?", labelCode).First(&label).Error; err != nil && err != gorm.ErrRecordNotFound {
		return label, err
	}
	if err != gorm.ErrRecordNotFound {
		label.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
		label.DeletedBy = userCode
		if err = config.MysqlDB.Save(&label).Error; err != nil && err != gorm.ErrRecordNotFound {
			return label, err
		}
	}
	// 从缓存中删除
	utils.RedisRemove(utils.UserLabelCachePrefix + labelCode)
	return label, nil
}

// GetNextNotAudited 获取下一个待审核的标签
func (labelService *LabelService) GetNextNotAudited() (db.UserLabelModel, error) {
	var label db.UserLabelModel
	err := config.MysqlDB.Where("status <= 6 and status >= 4").First(&label).Error
	return label, err
}
