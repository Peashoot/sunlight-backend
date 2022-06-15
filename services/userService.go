package services

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/peashoot/sunlight/auth"
	"github.com/peashoot/sunlight/config"
	"github.com/peashoot/sunlight/entity/db"
	"github.com/peashoot/sunlight/log"
	"github.com/peashoot/sunlight/utils"
	"gorm.io/gorm"
)

const (
	telUserRegAccountPrefix = "sunlight_" // 手机用户注册账号前缀
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

// GetUserByAccount 根据用户账号（或手机号）获取用户实体
func (userService *UserService) GetUserByAccount(account string) (db.UserModel, error) {
	// 先从缓存根据用户账号或手机号查询是否为无效的信息
	log.Info("[UserService.GetUserByAccount]", "begin check user", account, "from cache")
	var user db.UserModel
	_, err := utils.RedisGet[int64](utils.UnvalidUserCachePrefix + account)
	// 查询到无效用户账号
	if err == nil {
		log.Info("[UserService.GetUserByAccount]", "check unvalid user", account, "from cache")
		return user, err
	}
	// 缓存查询不到，到数据库查询
	log.Info("[UserService.GetUserByAccount]", "begin find user", account, "from mysql")
	reg := regexp.MustCompile(`^\d+$`)
	var result *gorm.DB
	// 如果是纯数字，那就是手机号；其他就用账号去查
	if reg.Match([]byte(account)) {
		result = config.MysqlDB.Where("telephone = ?", account).First(&user)
	} else {
		result = config.MysqlDB.Where("account = ?", account).First(&user)
	}
	if result.Error == gorm.ErrRecordNotFound { // 没查到数据
		// 添加一条无效记录到缓存
		utils.RedisSetT(utils.UnvalidUserCachePrefix+account, time.Now().UnixMilli(),
			time.Minute*time.Duration(config.GetValue[int](config.RCN_ValidUserCacheExpiration)))
		return user, utils.UserNotFoundError{Message: "cannot found user " + account}
	}
	if result.Error != nil { // mysql查询失败
		return user, result.Error
	}
	return user, nil
}

// GetUserByCode 根据用户编号获取用户实体
func (userService *UserService) GetUserByCode(code string) (db.UserModel, error) {
	// 先从缓存根据用户编号查询用户信息
	log.Info("[UserService.GetUserByCode]", "begin check user", code, "from cache")
	user, err := utils.RedisGet[db.UserModel](utils.UserRedisCachePrefix + code)
	// 查询到无效用户账号
	if err == nil {
		log.Info("[UserService.GetUserByCode]", "found user", code, "from cache")
		return user, err
	}
	// 缓存查询不到，到数据库查询
	log.Info("[UserService.GetUserByCode]", "begin find user", code, "from mysql")
	// 如果是纯数字，那就是手机号；其他就用账号去查
	result := config.MysqlDB.Where("code = ?", code).First(&user)
	if result.Error != nil { // mysql查询失败
		return user, result.Error
	}
	if result.RowsAffected == 0 { // 没查到数据
		return user, utils.UserNotFoundError{Message: "cannot found user " + code}
	}
	// 更新用户信息到缓存
	utils.RedisSetT(utils.UserRedisCachePrefix+user.Code, user,
		time.Minute*time.Duration(config.GetValue[int](config.RCN_UserCacheExpiration)))
	return user, nil
}

// UpdateProfile 更新用户信息
func (userService *UserService) UpdateProfile(user *db.UserModel) error {
	log.Info("[UserService.UpdateUser]", "update user", user.Code, "in db")
	user.UpdatedBy = user.Code
	user.UpdatedAt = time.Now()
	result := config.MysqlDB.Save(&user)
	if result.Error != nil {
		return result.Error
	}
	log.Info("[UserService.UpdateUser]", "update user", user.Code, "in cache")
	utils.RedisSetT(utils.UserRedisCachePrefix+user.Code, user,
		time.Minute*time.Duration(config.GetValue[int](config.RCN_UserCacheExpiration)))
	return nil
}

// UpdateKeyInfo 更新用户关键信息
func (userService *UserService) UpdateKeyInfo(user *db.UserModel, account, telephone, password string) error {
	// 清除用户缓存
	var count int64
	// 检查是否存在重复用户名
	if account != "" {
		if err := config.MysqlDB.Model(&db.UserModel{}).Where("account = ? and id != ?", account, user.ID).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return utils.RecordExistsFoundError{Message: "duplicate account"}
		}
		user.Account = account
	}
	// 检查是否存在重复手机号
	if telephone != "" {
		if err := config.MysqlDB.Model(&db.UserModel{}).Where("telephone = ? and id != ?", telephone, user.ID).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return utils.RecordExistsFoundError{Message: "duplicate telephone"}
		}
		user.Telephone = telephone
	}
	// 密码赋值
	if password != "" {
		user.Password = user.GetPwdHash(password)
	}
	log.Info("[UserService.UpdateKeyInfo]", "update user key info to db")
	user.UpdatedBy = user.Code
	if err := config.MysqlDB.Save(&user).Error; err != nil {
		return err
	}
	// 更新缓存
	utils.RedisSetT(utils.UserRedisCachePrefix+user.Code, &user,
		time.Minute*time.Duration(config.GetValue[int](config.RCN_UserCacheExpiration)))
	// 清除无效用户缓存
	utils.RedisRemove(utils.UnvalidUserCachePrefix + user.Account)
	utils.RedisRemove(utils.UnvalidUserCachePrefix + user.Telephone)
	return nil
}

// GenerateAuthToken 生成jwt
func (userService *UserService) GenerateAuthToken(user *db.UserModel) (string, error) {
	claims := auth.AuthClaims{
		Account:        user.Account,
		Code:           user.Code,
		LoginTimestamp: time.Now().Unix(),
	}
	tokenBytes, err := auth.JwtSigner.Sign(claims)
	if err != nil {
		return "", err
	}
	token := string(tokenBytes)
	// 新增或覆盖原来生成的Token，单点登录
	authTokenModel := db.AuthTokenRecord{
		BaseModel: db.BaseModel{
			Code: user.Code,
		},
	}
	// 把Token更新到数据库
	result := config.MysqlDB.Where(authTokenModel).FirstOrCreate(&authTokenModel).Update("token", token)
	if result.Error != nil {
		return "", result.Error
	} else if result.RowsAffected < 1 {
		return "", utils.DBExecuteError{Message: "fail to update token into DB"}
	}
	// 把Token更新到redis
	utils.RedisSetT(utils.JwtRedisCachePrefix+user.Code, token,
		time.Minute*time.Duration(config.GetValue[int](config.RCN_UserTokenCacheExpiration)))
	return token, nil
}

// ClearAuthToken 清除认证Token
func (userService *UserService) ClearAuthToken(userCode string) error {
	// 清除用户token缓存
	if err := utils.RedisRemove(userCode); err != nil {
		return err
	}
	// 清除数据库Token缓存
	result := config.MysqlDB.Model(&db.AuthTokenRecord{}).Where("code = ?", userCode).Update("token", "")
	return result.Error
}

// RegisterUser 注册用户
func (userService *UserService) RegisterUser(telephone, password string) (db.UserModel, error) {
	user := db.UserModel{
		BaseModel: db.BaseModel{
			Code: uuid.NewString(),
		},
		Account:   telUserRegAccountPrefix + time.Now().Format("060102") + hex.EncodeToString(md5.New().Sum([]byte(telephone)))[8:16] + fmt.Sprintf("%4d", rand.Int31n(10000)) + telephone[len(telephone)-5:],
		Telephone: telephone,
		Salt:      strings.ReplaceAll(uuid.NewString(), "-", "")[8:24],
		Birthday:  time.Now(),
	}
	user.NickName = user.Account
	user.Password = user.GetPwdHash(password)
	user.CreatedBy = user.Code
	// 先把用户添加到数据库中
	result := config.MysqlDB.Create(&user)
	if result.Error != nil {
		return user, result.Error
	}
	if result.RowsAffected < 1 {
		return user, utils.DBExecuteError{Message: "fail to insert user"}
	}
	// 更新缓存
	utils.RedisSetT(utils.UserRedisCachePrefix+user.Code, user,
		time.Minute*time.Duration(config.GetValue[int](config.RCN_UserCacheExpiration)))
	// 清除无效用户缓存
	utils.RedisRemove(utils.UnvalidUserCachePrefix + user.Account)
	utils.RedisRemove(utils.UnvalidUserCachePrefix + user.Telephone)
	return user, nil
}

// RemoveUser 注销账户
func (userService *UserService) RemoveUser(user *db.UserModel, opUserCode string) error {
	// 更新删除人，删除时间（软删除只能修改时间）
	log.Info("[UserService.RemoveUser]", "remove user", user.Code, "from db")
	result := config.MysqlDB.Model(user).Updates(
		db.UserModel{BaseModel: db.BaseModel{DeletedBy: opUserCode, DeletedAt: gorm.DeletedAt{Time: time.Now(), Valid: true}}})
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}
	// 清除redis用户缓存
	log.Info("[UserService.RemoveUser]", "clear user", user.Code, "info in cache")
	utils.RedisRemove(utils.UserRedisCachePrefix + user.Code)
	return nil
}
