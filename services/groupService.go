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

type UserGroupService struct {
}

func NewGroupService() *UserGroupService {
	return &UserGroupService{}
}

// CreateNewGroup 创建一个新的群组
func (groupService *UserGroupService) CreateNewGroup(groupName, opUserCode string, inviteeCodes ...string) (db.UserGroupModel, error) {
	// 创建群组信息
	group := db.UserGroupModel{
		BaseModel: db.BaseModel{
			Code:      uuid.NewString(),
			CreatedBy: opUserCode,
		},
		Name: groupName,
	}
	memberships := make([]db.GroupMembershipModel, len(inviteeCodes))
	for index, inviteeCode := range inviteeCodes {
		memberships[index] = db.GroupMembershipModel{
			BaseModel: db.BaseModel{
				Code:      uuid.NewString(),
				CreatedBy: opUserCode,
			},
			GroupCode:  group.Code,
			MemberCode: inviteeCode,
		}
	}
	// 将群组信息和群组成员保存到数据库
	if err := config.MysqlDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&group).Error; err != nil {
			return err
		}
		for _, membership := range memberships {
			if err := tx.Create(&membership).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return group, err
	}
	// 将群组信息和群组成员存到缓存中
	utils.RedisSetT(utils.UserGroupCachePrefix+group.Code, group,
		time.Minute*time.Duration(config.GetValue[int](config.RCN_GroupCacheExpiration)))
	for _, membership := range memberships {
		utils.RedisHSet(utils.GroupMembershipCachePrefix+group.Code, membership.Code, membership)
	}
	utils.RedisExpire(utils.GroupMembershipCachePrefix+group.Code,
		time.Minute*time.Duration(config.GetValue[int](config.RCN_GroupMemberCacheExpiration)))
	return group, nil
}

// InviteUser 邀请新成员加入当前群组
func (groupService *UserGroupService) InviteUser(groupCode, opUserCode string, inviteeCodes ...string) (db.UserGroupModel, error) {
	// 从缓存中查询群组信息
	log.Info("[UserGroupService.InviteUser]", "try to get group", groupCode, "from cache")
	group, err := utils.RedisGet[db.UserGroupModel](utils.UserGroupCachePrefix + groupCode)
	if err != nil {
		log.Info("[UserGroupService.InviteUser]", "find group", groupCode, "from db")
		if err = config.MysqlDB.Where("code = ?", groupCode).First(&group).Error; err != nil {
			return group, err
		}
	}
	log.Info("[UserGroupService.InviteUser]", "add new member to group")
	memberships := make([]db.GroupMembershipModel, len(inviteeCodes))
	for index, inviteeCode := range inviteeCodes {
		// 如果数据库已存在当前code的组员信息，停止添加
		var count int64
		err = config.MysqlDB.Model(&db.GroupMembershipModel{}).Where("code = ?", inviteeCode).Count(&count).Error
		if err != nil {
			return group, err
		}
		if count > 0 {
			return group, utils.RecordExistsFoundError{Message: "code " + inviteeCode + " already exists"}
		}
		memberships[index] = db.GroupMembershipModel{
			BaseModel: db.BaseModel{
				Code:      uuid.NewString(),
				CreatedBy: opUserCode,
			},
			GroupCode:  groupCode,
			MemberCode: inviteeCode,
		}
	}
	// 将群组成员保存到数据库
	if err := config.MysqlDB.Transaction(func(tx *gorm.DB) error {
		for _, membership := range memberships {
			if err := tx.Create(&membership).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return group, err
	}
	// 添加信息到缓存
	for _, membership := range memberships {
		utils.RedisHSet(utils.GroupMembershipCachePrefix+group.Code, membership.Code, membership)
	}
	utils.RedisExpire(utils.GroupMembershipCachePrefix+group.Code,
		time.Minute*time.Duration(config.GetValue[int](config.RCN_GroupMemberCacheExpiration)))
	return group, err
}

// ChangeGroupName 改变群组名称
func (groupService *UserGroupService) ChangeGroupName(groupCode, newGroupName, opUserCode string) (db.UserGroupModel, error) {
	log.Info("[UserGroupService.ChangeGroupName]", "get group", groupCode, "info")
	group, err := groupService.GetGroupModelByCode(groupCode)
	if err != nil {
		return group, err
	}
	// 修改群组信息
	log.Info("[UserGroupService.ChangeGroupName]", "modify group", groupCode, "info in db")
	group.Name = newGroupName
	group.UpdatedBy = opUserCode
	group.UpdatedAt = time.Now()
	if err = config.MysqlDB.Save(&group).Error; err != nil {
		return group, err
	}
	// 更新缓存
	log.Info("[UserGroupService.ChangeGroupName]", "write group", groupCode, "info into cache")
	utils.RedisSetT(utils.UserGroupCachePrefix+groupCode, group,
		time.Minute*time.Duration(config.GetValue[int](config.RCN_GroupCacheExpiration)))
	return group, nil
}

// ChangeMemberNickname 改变成员昵称
func (grouService *UserGroupService) ChangeMemberNickname(groupCode, memberCode, newNickname, opUserCode string) (db.GroupMembershipModel, error) {
	var membership db.GroupMembershipModel
	log.Info("[UserGroupService.ChangeMemberNickname]", "get membership", groupCode, memberCode, "from cache")
	exists, err := utils.RedisExists(utils.GroupMembershipCachePrefix + groupCode)
	if exists {
		membership, err = utils.RedisHGet[db.GroupMembershipModel](utils.GroupMembershipCachePrefix+groupCode, memberCode)
	}
	if err != nil || !exists {
		log.Info("[UserGroupService.ChangeMemberNickname]", "get membership", groupCode, memberCode, "from db")
		if err = config.MysqlDB.Where("group_code = ? and member_code = ?", groupCode, memberCode).First(&membership).Error; err != nil {
			return membership, err
		}
	}
	membership.MemberName = newNickname
	membership.UpdatedBy = opUserCode
	membership.UpdatedAt = time.Now()
	if err = config.MysqlDB.Save(&membership).Error; err != nil {
		return membership, err
	}
	// 如果缓存中存在，刷新缓存中的值；不存在就不刷新，下次从数据库中取
	log.Info("[UserGroupService.QuitGroup]", "try to update membership", groupCode, memberCode, "in cache")
	exists, _ = utils.RedisExists(utils.GroupMembershipCachePrefix + groupCode)
	if exists {
		utils.RedisHSet(utils.GroupMembershipCachePrefix+groupCode, memberCode, &membership)
	}
	return membership, nil
}

// GetGroupModelByCode 根据编号获取群组信息
func (groupService *UserGroupService) GetGroupModelByCode(groupCode string) (db.UserGroupModel, error) {
	// 从缓存中获取
	group, err := utils.RedisGet[db.UserGroupModel](utils.UserGroupCachePrefix + groupCode)
	if err == nil {
		return group, nil
	}
	// 从数据库获取
	err = config.MysqlDB.Where("code = ?", groupCode).First(&group).Error
	return group, err
}

// EraseGroup 解散群组信息（包括所有成员）
func (groupService *UserGroupService) EraseGroup(groupCode, opUserCode string) (db.UserGroupModel, error) {
	var group db.UserGroupModel
	// 从数据库删除数据
	log.Info("[UserGroupService.EraseGroup]", "remove group", groupCode, "info and membership from db")
	if err := config.MysqlDB.Transaction(func(tx *gorm.DB) error {
		// 删除所有群组成员
		if err := tx.Model(&db.GroupMembershipModel{}).Where("group_code = ?", groupCode).Updates(
			&db.GroupMembershipModel{BaseModel: db.BaseModel{DeletedBy: opUserCode,
				DeletedAt: gorm.DeletedAt{Time: time.Now(), Valid: true}}}).Error; err != nil {
			return err
		}
		// 删除群组信息
		if err := tx.Model(&db.UserGroupModel{}).Where("code = ?", groupCode).First(&group).Error; err != nil {
			return err
		}
		group.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
		group.DeletedBy = opUserCode
		if err := tx.Save(group).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return group, err
	}
	// 同步删除缓存
	log.Info("[UserGroupService.EraseGroup]", "remove group", groupCode, "from cache")
	utils.RedisRemove(utils.GroupMembershipCachePrefix+groupCode, utils.UserGroupCachePrefix+groupCode)
	return group, nil
}

// QuitGroup 退出群组
func (groupService *UserGroupService) QuitGroup(groupCode, quitedCode, opUserCode string) (db.GroupMembershipModel, error) {
	var membership db.GroupMembershipModel
	log.Info("[UserGroupService.QuitGroup]", "get membership", groupCode, quitedCode, "from cache")
	exists, err := utils.RedisExists(utils.GroupMembershipCachePrefix + groupCode)
	if exists {
		membership, err = utils.RedisHGet[db.GroupMembershipModel](utils.GroupMembershipCachePrefix+groupCode, quitedCode)
	}
	if err != nil || !exists {
		log.Info("[UserGroupService.QuitGroup]", "get membership", groupCode, quitedCode, "from db")
		if err = config.MysqlDB.Where("group_code = ? and member_code = ?", groupCode, quitedCode).First(&membership).Error; err != nil {
			return membership, err
		}
	}
	membership.DeletedBy = opUserCode
	membership.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
	if err = config.MysqlDB.Save(&membership).Error; err != nil {
		return membership, err
	}
	// 如果缓存中存在，刷新缓存中的值；不存在就不刷新，下次从数据库中取
	log.Info("[UserGroupService.QuitGroup]", "try to remove membership", groupCode, quitedCode, "from cache")
	exists, _ = utils.RedisExists(utils.GroupMembershipCachePrefix + groupCode)
	if exists {
		utils.RedisHDel(utils.GroupMembershipCachePrefix+groupCode, quitedCode)
	}
	return membership, nil
}

// QueryGroupByMemberCode 用户查询所有群组
func (groupService *UserGroupService) QueryGroupByMemberCode(memberCode string) ([]db.UserGroupModel, error) {
	groups := make([]db.UserGroupModel, 0)
	// 从数据库查询用户群组关系
	log.Info("[UserGroupService.QueryGroupByMemberCode]", "find groupship of user", memberCode, "from db")
	var memships []db.GroupMembershipModel
	if err := config.MysqlDB.Model(&db.GroupMembershipModel{}).Where("member_code = ?", memberCode).Find(&memships).Error; err != nil {
		return groups, err
	}
	for _, memship := range memships {
		group, err := groupService.GetGroupModelByCode(memship.GroupCode)
		if err == nil {
			groups = append(groups, group)
		}
	}
	return groups, nil
}
