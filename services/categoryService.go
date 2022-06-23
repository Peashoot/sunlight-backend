package services

import (
	"time"

	"github.com/google/uuid"
	"github.com/peashoot/sunlight/config"
	"github.com/peashoot/sunlight/entity/db"
	"github.com/peashoot/sunlight/log"
	"github.com/peashoot/sunlight/utils"
	"gorm.io/gorm"
)

type CategoryService struct {
}

func NewCategoryService() *CategoryService {
	return &CategoryService{}
}

// Create 创建新类别
func (categoryService *CategoryService) Create(name, groupCode, opUserCode, desc string) (db.CategoryModel, error) {
	log.Info("[CategoryService.Create]", "create a model an save into db")
	category := db.CategoryModel{
		BaseModel: db.BaseModel{
			Code:      uuid.NewString(),
			CreatedBy: opUserCode,
		},
		Name:        name,
		OwnerGroup:  groupCode,
		Description: desc,
	}
	if err := config.MysqlDB.Create(&category).Error; err != nil {
		return category, err
	}
	log.Info("[CategoryService.Create]", "add model into cache")
	utils.RedisSetT(utils.FileCategoryCachePrefix+category.Code, &category,
		time.Minute*time.Duration(config.GetValue[int](config.RCN_CategoryCacheExpiration)))
	return category, nil
}

// ChangeName 修改类别的名称
func (categoryService *CategoryService) ChangeName(name, categoryCode, opUserCode string) (db.CategoryModel, error) {
	var category db.CategoryModel
	log.Info("[CategoryService.ChangeName]", "update name of category", categoryCode, "in db")
	if err := config.MysqlDB.Where("code = ?", categoryCode).First(&category).Error; err != nil {
		return category, err
	}
	category.Name = name
	category.UpdatedBy = opUserCode
	category.UpdatedAt = time.Now()
	if err := config.MysqlDB.Save(&category).Error; err != nil {
		return category, nil
	}
	log.Info("[CategoryService.ChangeName]", "update category", categoryCode, "in cache")
	utils.RedisSetT(utils.FileCategoryCachePrefix+categoryCode, &category,
		time.Minute*time.Duration(config.GetValue[int](config.RCN_CategoryCacheExpiration)))
	return category, nil
}

// ChangeDesc 修改类别的描述
func (categoryService *CategoryService) ChangeDesc(desc, categoryCode, opUserCode string) (db.CategoryModel, error) {
	var category db.CategoryModel
	log.Info("[CategoryService.ChangeDesc]", "update description of category", categoryCode, "in db")
	if err := config.MysqlDB.Where("code = ?", categoryCode).First(&category).Error; err != nil {
		return category, err
	}
	category.Description = desc
	category.UpdatedBy = opUserCode
	category.UpdatedAt = time.Now()
	if err := config.MysqlDB.Save(&category).Error; err != nil {
		return category, nil
	}
	log.Info("[CategoryService.ChangeDesc]", "update category", categoryCode, "in cache")
	utils.RedisSetT(utils.FileCategoryCachePrefix+categoryCode, &category,
		time.Minute*time.Duration(config.GetValue[int](config.RCN_CategoryCacheExpiration)))
	return category, nil
}

// RemoveByGroupCode 根据群组编号删除
func (categoryService *CategoryService) RemoveByGroupCode(groupCode, opUserCode string) ([]db.CategoryModel, error) {
	// 清除数据库
	var categories []db.CategoryModel
	config.MysqlDB.Transaction(func(tx *gorm.DB) error {
		// 根据群组编号查询出所有记录
		log.Info("[CategoryService.Remove]", "remove category from cache")
		if err := tx.Where("group_code = ?", groupCode).Find(&categories).Error; err != nil {
			return err
		}
		for _, category := range categories {
			// 删除记录
			category.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
			category.DeletedBy = opUserCode
			if err := tx.Save(&category).Error; err != nil {
				return err
			}
			// 把对应文件上的分组信息也删除
			if err := tx.Model(&db.FileStorageRecord{}).Where("category_code = ?", category.Code).Updates(&db.FileStorageRecord{
				BaseModel: db.BaseModel{
					UpdatedBy: opUserCode,
				},
				CategoryCode: "",
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
	// 清除缓存
	log.Info("[CategoryService.RemoveByGroupCode]", "remove category info from cache")
	for _, category := range categories {
		utils.RedisRemove(utils.FileCategoryCachePrefix + category.Code)
	}
	return categories, nil
}

// Remove 移除类别
func (categoryService *CategoryService) Remove(categoryCode, opUserCode string) (db.CategoryModel, error) {
	// 删除数据库中类别实体
	var category db.CategoryModel
	log.Info("[CategoryService.Remove]", "remove category", categoryCode, "from cache")
	if err := config.MysqlDB.Where("code = ?", categoryCode).First(&category).Error; err != nil {
		return category, err
	}
	category.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
	category.DeletedBy = opUserCode
	if err := config.MysqlDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&category).Error; err != nil {
			return err
		}
		// 清除文件记录中该类别信息
		log.Info("[CategoryService.Remove]", "wipe out category of files", categoryCode, "in db")
		if err := tx.Model(&db.FileStorageRecord{}).Where("category_code = ?", categoryCode).Updates(&db.FileStorageRecord{
			BaseModel: db.BaseModel{
				UpdatedBy: opUserCode,
			},
			CategoryCode: "",
		}).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return category, err
	}
	// 清除文件记录中该类别信息
	log.Info("[CategoryService.Remove]", "wipe out category of files", categoryCode, "in db")
	if err := config.MysqlDB.Model(&db.FileStorageRecord{}).Where("category_code = ?", categoryCode).Updates(&db.FileStorageRecord{
		BaseModel: db.BaseModel{
			UpdatedBy: opUserCode,
		},
		CategoryCode: "",
	}).Error; err != nil {
		return category, nil
	}
	// 从缓存中删除
	log.Info("[CategoryService.Remove]", "remove category", categoryCode, "from cache")
	utils.RedisRemove(utils.FileCategoryCachePrefix + categoryCode)
	return category, nil
}
