package db

import (
	"crypto/sha256"
	"fmt"
	"time"
)

// UserModel 用户信息
type UserModel struct {
	BaseModel
	Account       string     `gorm:"size:50;comment:账号;index"`   // 用户名
	NickName      string     `gorm:"size:50;comment:昵称;index"`   // 昵称
	Password      string     `gorm:"size:100;comment:密码"`        // 密码
	Salt          string     `gorm:"size:32;comment:盐"`          // 密码加盐
	RoleID        uint       `gorm:"comment:角色ID"`               // 角色ID
	Telephone     string     `gorm:"size:50;comment:手机号"`        // 手机号
	Email         string     `gorm:"size:50;comment:邮箱"`         // 邮箱
	Avatar        string     `gorm:"size:256;comment:头像链接"`      // 用户头像链接地址
	LastLoginIP   string     `gorm:"size:50;comment:上一次登录的IP地址"` // 上一次登录的IP地址
	LastLoginTime *time.Time `gorm:"comment:上一次登录的时间"`           // 上一次登录的时间
	Token         string     `gorm:"size:512;comment:登录授权校验"`    // 登录授权校验
	Birthday      time.Time  `gorm:"comment:用户生日"`               // 用户生日
	Bio           string     `gorm:"comment:个性签名"`               // 个性签名
	Labels        string     `gorm:"size:512;comment:个人标签"`      // 个人标签
}

// GetPwdHash 获取密码hash值
func (user *UserModel) GetPwdHash(password string) string {
	origin := []byte(user.Salt + password + user.Salt)
	hashBytes := sha256.Sum256(origin)
	return fmt.Sprintf("%Xn", hashBytes)
}
