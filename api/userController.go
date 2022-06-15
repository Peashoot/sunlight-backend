package api

import (
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"github.com/peashoot/sunlight/auth"
	"github.com/peashoot/sunlight/config"
	"github.com/peashoot/sunlight/entity/db"
	"github.com/peashoot/sunlight/entity/do/requests"
	"github.com/peashoot/sunlight/entity/do/responses"
	"github.com/peashoot/sunlight/log"
	"github.com/peashoot/sunlight/services"
	"github.com/peashoot/sunlight/utils"
)

// UserController 用户控制器
type UserController struct {
	userService *services.UserService
}

// NewUserController 创建一个用户控制器
func NewUserController() *UserController {
	return &UserController{
		userService: services.NewUserService(),
	}
}

// Route 路由注册
func (controller *UserController) Route(app iris.Party) {
	app.Post("/login", controller.Login)
	app.Post("/user/register", controller.Register)
	auth.AuthNeedParty.Post("/logout", controller.Logout)
	auth.AuthNeedParty.Post("/user/editprofile", controller.EditProfile)
	auth.AuthNeedParty.Post("/user/changekey", controller.ChangeKeyinfo)
}

// Login 用户登录 Post /login
// 分布式锁，同一时间内同一账号的请求，只处理一次，其他直接返回“系统繁忙”
// 登录失败的账号，把账号信息缓存起来，下次先检查是否为无效信息，是的话，直接返回失败
// 到数据库中查询出账号信息，生成Token，并返回基本信息
func (controller *UserController) Login(ctx iris.Context) {
	backDo := responses.NewPackagedRespModel()
	var locked bool
	var err error
	var loginDo requests.UserLoginReqModel
	defer func() {
		if err := recover(); err != nil {
			log.Error("[UserController.Login]", "fatal when try to login, the exception is", err)
			backDo.Code = responses.ErrorCode
			backDo.Msg = "系统异常"
		}
		if locked {
			// 解锁
			utils.RedisTryUnlock(utils.UserLoginLockPrefix + loginDo.Account)
		}
		ctx.JSON(backDo)
	}()
	if err := ctx.ReadJSON(&loginDo); err != nil {
		log.Error("[UserController.Login]", "try to read json appear an exception:", err)
		return
	}
	// 获取到锁，才执行
	locked, err = utils.RedisTryLock(utils.UserLoginLockPrefix+loginDo.Account,
		time.Duration(config.GetValue[int](config.RCN_UserLoginLockExpiration))*time.Second)
	if err != nil {
		log.Error("[UserController.Login]", "try to set login lock for user", loginDo.Account, ", appear an exception:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	if !locked {
		log.Error("[UserController.Login]", "user", loginDo.Account, "request multiple times in a short time")
		backDo.Code = responses.SystemBusyCode
		backDo.Msg = "系统繁忙，请稍后重试"
		return
	}
	// 查询用户信息
	user, err := controller.userService.GetUserByAccount(loginDo.Account)
	if err != nil {
		if _, ok := err.(utils.UserNotFoundError); ok {
			log.Error("[UserController.Login]", "cannot find user", loginDo.Account, ", appear an exception:", err)
			backDo.Code = responses.NotFoundCode
			backDo.Msg = "用户名或密码错误"
			return
		}
		log.Error("[UserController.Login]", "cannot find user", loginDo.Account, ", appear an exception:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	// 校验密码
	if user.GetPwdHash(loginDo.Password) != user.Password {
		log.Error("[UserController.Login]", "password of user", loginDo.Account, "input wrong")
		backDo.Code = responses.NotFoundCode
		backDo.Msg = "用户名或密码错误"
		return
	}
	// 生成jwt
	token, err := controller.userService.GenerateAuthToken(&user)
	if err != nil {
		log.Error("[UserController.Login]", "generate token for user", user.Account, "appear an exception:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	// 添加操作记录
	realIP := ctx.Request().Header.Get("X-Real-IP")
	go services.AddOperateRecord(user.Code, user.Code, db.UserDataType, db.LoadActionType,
		"user ", user.Code, " log in at", realIP)
	backDo.Code = responses.OKCode
	backDo.Msg = "成功"
	respModel := &responses.UserLoginRespModel{
		Token:    token,
		Avatar:   user.Avatar,
		NickName: user.NickName,
		UserCode: user.Code,
	}
	backDo.Data = respModel
}

// Logout 用户登出 Post /logout
// 清除用户Token信息
func (controller *UserController) Logout(ctx iris.Context) {
	backDo := responses.NewPackagedRespModel()
	defer func() {
		if err := recover(); err != nil {
			log.Error("[UserController.Logout]", "fatal when try to log out, the exception is", err)
			backDo.Code = responses.ErrorCode
			backDo.Msg = "系统异常"
		}
		ctx.JSON(backDo)
	}()
	authClaims := jwt.Get(ctx).(*auth.AuthClaims)
	if err := controller.userService.ClearAuthToken(authClaims.Code); err != nil {
		log.Error("[UserController.Logout]", "fail to clear", authClaims.Code, "token, appear an exception:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "操作失败"
		return
	}
	ctx.Logout()
	// 添加操作记录
	go services.AddOperateRecord(authClaims.Code, authClaims.Code, db.UserDataType, db.CloseActionType,
		"user ", authClaims.Code, " log out")
	backDo.Code = responses.OKCode
	backDo.Msg = "成功"
}

// Register 用户注册 Post /user/register
// 分布式锁，同一时间的同一手机号的注册，晚到的请求不再处理
func (controller *UserController) Register(ctx iris.Context) {
	backDo := responses.NewPackagedRespModel()
	locked := false
	var err error
	var registerDo requests.UserRegisterReqModel
	defer func() {
		if err := recover(); err != nil {
			log.Error("[UserController.Register]", "fatal when try to register, the exception is", err)
			backDo.Code = responses.ErrorCode
			backDo.Msg = "系统异常"
		}
		if locked {
			// 解锁
			utils.RedisTryUnlock(utils.UserRegLockPrefix + registerDo.Telephone)
		}
		ctx.JSON(backDo)
	}()
	if err = ctx.ReadJSON(&registerDo); err != nil {
		log.Error("[UserController.Register]", "try to read json appear an exception:", err)
		return
	}
	// 参数校验
	if registerDo.Telephone == "" {
		log.Info("[UserController.Register]", "empty telephone is not allowed when register")
		backDo.Code = responses.NotAllowedCode
		backDo.Msg = "手机号不能为空"
		return
	}
	// 获取到锁，才执行
	locked, err = utils.RedisTryLock(utils.UserRegLockPrefix+registerDo.Telephone, time.Duration(config.GetValue[int](config.RCN_UserRegisterLockExpiration))*time.Second)
	if err != nil {
		log.Error("[UserController.Register]", "try to set reg lock for user", registerDo.Telephone, ", appear an exception:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	if !locked {
		log.Error("[UserController.Register]", "user", registerDo.Telephone, "request multiple times in a short time")
		backDo.Code = responses.SystemBusyCode
		backDo.Msg = "系统繁忙，请稍后重试"
		return
	}
	// 检查是否存在相同手机号的用户
	user, err := controller.userService.GetUserByAccount(registerDo.Telephone)
	if err != nil {
		if _, ok := err.(utils.UserNotFoundError); !ok {
			log.Error("[UserController.Register]", "try to find user", registerDo.Telephone, ", appear an exception:", err)
			backDo.Code = responses.ErrorCode
			backDo.Msg = "系统异常"
			return
		}
	}
	// 如果用户存在，返回失败
	if user.ID > 0 {
		log.Info("[UserController.Register]", "user", registerDo.Telephone, "has exists")
		backDo.Code = responses.BadRequestCode
		backDo.Msg = "用户已存在"
		return
	}
	// 注册用户
	if user, err = controller.userService.RegisterUser(registerDo.Telephone, registerDo.Password); err != nil {
		log.Error("[UserController.Register]", "try to register user", registerDo.Telephone, ", appear an exception:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	// 添加操作记录
	realIP := ctx.Request().Header.Get("X-Real-IP")
	go services.AddOperateRecord(user.Code, user.Code, db.UserDataType, db.InsertActionType,
		"user ", user.Code, " register in at ", realIP)
	backDo.Code = responses.CreatedCode
	backDo.Msg = "成功"
}

// EditProfile 修改用户信息 Post /user/editprofile
func (controller *UserController) EditProfile(ctx iris.Context) {
	backDo := responses.NewPackagedRespModel()
	defer func() {
		if err := recover(); err != nil {
			log.Error("[UserController.EditProfile]", "fatal when try to edit profile, the exception is", err)
			backDo.Code = responses.ErrorCode
			backDo.Msg = "系统异常"
		}
		ctx.JSON(backDo)
	}()
	var editDo requests.UserEditProfileReqModel
	if err := ctx.ReadJSON(&editDo); err != nil {
		log.Error("[UserController.EditProfile]", "try to read json, appear an exception:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	user, err := controller.userService.GetUserByCode(editDo.UserCode)
	if err != nil {
		if _, ok := err.(utils.UserNotFoundError); ok {
			log.Info("[UserController.EditProfile]", "cannot found user", editDo.UserCode)
			backDo.Code = responses.NotFoundCode
			backDo.Msg = "未找到用户信息"
			return
		}
		log.Error("[UserController.EditProfile]", "fail to get user", editDo.UserCode, ", appear an exception:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	editDo.EditUserAssignment(&user)
	if err := controller.userService.UpdateProfile(&user); err != nil {
		log.Error("[UserController.EditProfile]", "fail to update user", editDo.UserCode, ", appear an exception:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	// 添加操作记录
	realIP := ctx.Request().Header.Get("X-Real-IP")
	go services.AddOperateRecord(user.Code, user.Code, db.UserDataType, db.UpdateActionType,
		"user ", user.Code, " edit profile at ", realIP)
	backDo.Code = responses.OKCode
	backDo.Msg = "成功"
}

// ChangeKeyinfo 修改关键信息 /user/changekey
// 分布式锁，同一时间同一用户只允许修改一次
func (controller *UserController) ChangeKeyinfo(ctx iris.Context) {
	backDo := responses.NewPackagedRespModel()
	locked := false
	var err error
	var changeKeyInfoDo requests.UserChangeKeyInfoReqModel
	defer func() {
		if err := recover(); err != nil {
			log.Error("[UserController.ChangeKeyinfo]", "fatal when try to change key info, the exception is", err)
			backDo.Code = responses.ErrorCode
			backDo.Msg = "系统异常"
		}
		if locked {
			// 解锁
			utils.RedisTryUnlock(utils.UserCagKeyLockPrefix + changeKeyInfoDo.UserCode)
		}
		ctx.JSON(backDo)
	}()
	if err = ctx.ReadJSON(&changeKeyInfoDo); err != nil {
		log.Error("[UserController.ChangeKeyinfo]", "try to read json, appear an exception:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	// 获取到锁，才执行
	locked, err = utils.RedisTryLock(utils.UserCagKeyLockPrefix+changeKeyInfoDo.UserCode,
		time.Duration(config.GetValue[int](config.RCN_UserModifyLockExpiration))*time.Second)
	if err != nil {
		log.Error("[UserController.ChangeKeyinfo]", "try to set reg lock for user", changeKeyInfoDo.UserCode, ", appear an exception:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	if !locked {
		log.Error("[UserController.ChangeKeyinfo]", "user", changeKeyInfoDo.UserCode, "request multiple times in a short time")
		backDo.Code = responses.SystemBusyCode
		backDo.Msg = "系统繁忙，请稍后重试"
		return
	}
	user, err := controller.userService.GetUserByCode(changeKeyInfoDo.UserCode)
	if err != nil {
		if _, ok := err.(utils.UserNotFoundError); ok {
			log.Info("[UserController.ChangeKeyinfo]", "cannot found user", changeKeyInfoDo.UserCode)
			backDo.Code = responses.NotFoundCode
			backDo.Msg = "未找到用户信息"
			return
		}
		log.Error("[UserController.ChangeKeyinfo]", "fail to get user", changeKeyInfoDo.UserCode, ", appear an exception:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	if err := controller.userService.UpdateKeyInfo(&user, changeKeyInfoDo.NewAccount,
		changeKeyInfoDo.NewTelephone, changeKeyInfoDo.NewPassword); err != nil {
		if _, ok := err.(utils.RecordExistsFoundError); ok {
			log.Error("[UserController.ChangeKeyinfo]", "duplicate info user", changeKeyInfoDo.UserCode)
			backDo.Code = responses.BadRequestCode
			backDo.Msg = "信息已被使用"
			return
		}
		log.Error("[UserController.ChangeKeyinfo]", "fail to update user", changeKeyInfoDo.UserCode, ", appear an exception:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	// 添加操作记录
	realIP := ctx.Request().Header.Get("X-Real-IP")
	go services.AddOperateRecord(user.Code, user.Code, db.UserDataType, db.UpdateActionType,
		"user ", user.Code, " change key info at ", realIP)
	backDo.Code = responses.OKCode
	backDo.Msg = "成功"
}

// CancelAccount 注销账户 /user/canceluser
// 分布式锁，同一时间同一用户只允许注销一次
func (controller *UserController) CancelAccount(ctx iris.Context) {
	backDo := responses.NewPackagedRespModel()
	locked := false
	var cancelDo requests.UserCancelAccountReqModel
	var err error
	defer func() {
		if err := recover(); err != nil {
			log.Error("[UserController.CancelAccount]", "fatal when try to remove user, the exception is", err)
			backDo.Code = responses.ErrorCode
			backDo.Msg = "系统异常"
		}
		if locked {
			// 解锁
			utils.RedisTryUnlock(utils.UserDelKeyLockPrefix + cancelDo.UserCode)
		}
		ctx.JSON(backDo)
	}()
	if err = ctx.ReadJSON(&cancelDo); err != nil {
		log.Error("[UserController.CancelAccount]", "try to read json, appear an exception:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	// 获取到锁，才执行
	locked, err = utils.RedisTryLock(utils.UserDelKeyLockPrefix+cancelDo.UserCode,
		time.Duration(config.GetValue[int](config.RCN_UserCancelLockExpiration))*time.Second)
	if err != nil {
		log.Error("[UserController.CancelAccount]", "try to set remove lock for user", cancelDo.UserCode, ", appear an exception:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	if !locked {
		log.Error("[UserController.CancelAccount]", "user", cancelDo.UserCode, "request multiple times in a short time")
		backDo.Code = responses.SystemBusyCode
		backDo.Msg = "系统繁忙，请稍后重试"
		return
	}
	user, err := controller.userService.GetUserByCode(cancelDo.UserCode)
	if err != nil {
		if _, ok := err.(utils.UserNotFoundError); ok {
			log.Info("[UserController.CancelAccount]", "cannot found user", cancelDo.UserCode)
			backDo.Code = responses.NotFoundCode
			backDo.Msg = "未找到用户信息"
			return
		}
		log.Error("[UserController.CancelAccount]", "fail to get user", cancelDo.UserCode, ", appear an exception:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	// 清除Token
	authClaims := jwt.Get(ctx).(*auth.AuthClaims)
	if err := controller.userService.ClearAuthToken(user.Code); err != nil {
		log.Error("[UserController.Logout]", "fail to clear", user.Code, "token, appear an exception:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "操作失败"
		return
	}
	ctx.Logout()
	// 移除用户
	if err := controller.userService.RemoveUser(&user, authClaims.Code); err != nil {
		log.Error("[UserController.CancelAccount]", "try to remove user", user.Code, ", appear an exception:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	// 添加操作记录
	realIP := ctx.Request().Header.Get("X-Real-IP")
	go services.AddOperateRecord(user.Code, authClaims.Code, db.UserDataType, db.UpdateActionType,
		"user ", authClaims.Code, " cancel user ", user.Code, " info at ", realIP)
	backDo.Code = responses.OKCode
	backDo.Msg = "成功"
}
