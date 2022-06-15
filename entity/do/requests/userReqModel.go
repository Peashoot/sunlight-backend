package requests

import (
	"time"

	"github.com/peashoot/sunlight/entity/db"
)

// UserLoginReqModel 用户登录实体
type UserLoginReqModel struct {
	BaseReqModel
	Account  string `json:"account"`  // 账号
	Password string `json:"password"` // 密码
}

// UserRegisterReqModel 用户注册实体
type UserRegisterReqModel struct {
	BaseReqModel
	Telephone string `json:"telephone"` // 手机号
	Password  string `json:"password"`  // 密码
}

// UserEditProfileReqModel 用户修改信息实体
type UserEditProfileReqModel struct {
	BaseReqModel
	UserCode string `json:"usercode"` // 用户编码
	NickName string `json:"nickname"` // 用户昵称
	Email    string `json:"email"`    // 邮箱
	Avatar   string `json:"avatar"`   // 头像链接
	Birthday int64  `json:"birthday"` // 年月日
	Bio      string `json:"bio"`      // 个性签名
	Labels   string `json:"labels"`   // 个人标签
}

// EditUserAssignment 传输实体给数据库实体赋值
func (editDo *UserEditProfileReqModel) EditUserAssignment(user *db.UserModel) {
	user.Avatar = editDo.Avatar
	user.Bio = editDo.Bio
	user.Birthday = time.Unix(editDo.Birthday, 0)
	user.Labels = editDo.Labels
	user.NickName = editDo.NickName
	user.Email = editDo.Email
}

// UserChangeKeyInfoReqModel 用户修改密码实体
type UserChangeKeyInfoReqModel struct {
	BaseReqModel
	UserCode     string `json:"usercode"` // 用户编码
	NewPassword  string `json:"newpwd"`   // 新密码
	NewAccount   string `json:"newacc"`   // 新账号
	NewTelephone string `json:"newtel"`   // 新手机号
}

// UserCancelAccountReqModel 用户注销实体
type UserCancelAccountReqModel struct {
	BaseReqModel
	UserCode string `json:"usercode"` // 用户编码
}
