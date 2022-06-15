package responses

type UserLoginRespModel struct {
	Token    string `json:"token"`    // 授权Token
	Avatar   string `json:"avatar"`   // 头像链接
	NickName string `json:"nickname"` // 用户昵称
	UserCode string `json:"usercode"` // 用户编号
}
