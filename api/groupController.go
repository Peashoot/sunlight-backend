package api

import (
	"github.com/kataras/iris/v12"
	"github.com/peashoot/sunlight/auth"
	"github.com/peashoot/sunlight/entity/db"
	"github.com/peashoot/sunlight/entity/do/requests"
	"github.com/peashoot/sunlight/entity/do/responses"
	"github.com/peashoot/sunlight/log"
	"github.com/peashoot/sunlight/services"
)

type GroupController struct {
	groupService *services.UserGroupService
}

func NewGroupController() *GroupController {
	return &GroupController{
		groupService: services.NewGroupService(),
	}
}

func (controller *GroupController) Route(app iris.Party) {
	auth.AuthNeedParty.Post("/group/create", controller.CreateNew)
	auth.AuthNeedParty.Post("/group/changeInfo", controller.ChangeInfo)
	auth.AuthNeedParty.Post("/group/invite", controller.InviteMember)
	auth.AuthNeedParty.Post("/group/kickout", controller.KickOutMember)
	auth.AuthNeedParty.Post("/group/customize", controller.ChangeMember)
	auth.AuthNeedParty.Post("/group/dismiss", controller.Dismiss)
}

// CreateNew 创建 /group/create
func (controller *GroupController) CreateNew(ctx iris.Context) {
	backDo := responses.NewPackagedRespModel()
	defer func() {
		if err := recover(); err != nil {
			log.Error("[GroupController.CreateNew]", "fatal when try to create, the exception is", err)
			backDo.Code = responses.ErrorCode
			backDo.Msg = "系统异常"
		}
		ctx.JSON(backDo)
	}()
	var createDo requests.GroupCreateReqModel
	if err := ctx.ReadJSON(&createDo); err != nil {
		log.Error("[GroupController.CreateNew]", "read json appear an exception, err:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	group, err := controller.groupService.CreateNewGroup(createDo.Name, createDo.Description, createDo.OperatorCode, createDo.Members...)
	if err != nil {
		log.Error("[GroupController.CreateNew]", "try to add new group", createDo.Name, ", appear an exception:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	go services.AddOperateRecord(group.Code, createDo.OperatorCode, db.GroupDataType, db.InsertActionType,
		"user", createDo.OperatorCode, "create a group", group.Code, "info:", group.Name)
	backDo.Code = responses.OKCode
	backDo.Msg = "成功"
	backDo.Data = &responses.GroupCreateRespModel{Code: group.Code}
}

// ChangeInfo 修改群组信息 /group/changeInfo
func (controller *GroupController) ChangeInfo(ctx iris.Context) {
	backDo := responses.NewPackagedRespModel()
	defer func() {
		if err := recover(); err != nil {
			log.Error("[GroupController.ChangeInfo]", "fatal when try to change info of group, the exception is", err)
			backDo.Code = responses.ErrorCode
			backDo.Msg = "系统异常"
		}
		ctx.JSON(backDo)
	}()
	var changeDo requests.GroupChangeReqModel
	if err := ctx.ReadJSON(&changeDo); err != nil {
		log.Error("[GroupController.ChangeInfo]", "read json appear an exception, err:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	group, err := controller.groupService.ChangeGroupInfo(changeDo.Code, changeDo.Name, changeDo.Description, changeDo.OwnerCode, changeDo.OperatorCode)
	if err != nil {
		log.Error("[GroupController.CreateNew]", "try to change info of group", changeDo.Name, ", appear an exception:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	go services.AddOperateRecord(group.Code, changeDo.OperatorCode, db.GroupDataType, db.InsertActionType,
		"user ", changeDo.OperatorCode, " change group ", group.Code, "'s info: ", group.Name)
	backDo.Code = responses.OKCode
	backDo.Msg = "成功"
}

// InviteMember 邀请成员加入群组 /group/invite
func (controller *GroupController) InviteMember(ctx iris.Context) {
	backDo := responses.NewPackagedRespModel()
	defer func() {
		if err := recover(); err != nil {
			log.Error("[GroupController.InviteMember]", "fatal when try to invite member into group, the exception is", err)
			backDo.Code = responses.ErrorCode
			backDo.Msg = "系统异常"
		}
		ctx.JSON(backDo)
	}()
	var inviteDo requests.GroupInviteMemberReqModel
	if err := ctx.ReadJSON(&inviteDo); err != nil {
		log.Error("[GroupController.InviteMember]", "read json appear an exception, err:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	group, err := controller.groupService.InviteUser(inviteDo.Code, inviteDo.OperatorCode, inviteDo.MemberCodes...)
	if err != nil {
		log.Error("[GroupController.InviteMember]", "try to invite member into group", inviteDo.Code, ", appear an exception:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	go services.AddOperateRecord(group.Code, inviteDo.OperatorCode, db.GroupMemberDataType, db.InsertActionType,
		"user ", inviteDo.OperatorCode, " invite member into group ", group.Code)
	backDo.Code = responses.OKCode
	backDo.Msg = "成功"
}

// KickOutMember 将成员踢出群组 /group/kickout
func (controller *GroupController) KickOutMember(ctx iris.Context) {
	backDo := responses.NewPackagedRespModel()
	defer func() {
		if err := recover(); err != nil {
			log.Error("[GroupController.KickOutMember]", "fatal when try to invite member into group, the exception is", err)
			backDo.Code = responses.ErrorCode
			backDo.Msg = "系统异常"
		}
		ctx.JSON(backDo)
	}()
	var kickoutDo requests.GroupKickoutMemberReqModel
	if err := ctx.ReadJSON(&kickoutDo); err != nil {
		log.Error("[GroupController.KickOutMember]", "read json appear an exception, err:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	group, err := controller.groupService.QuitGroup(kickoutDo.Code, kickoutDo.MemberCode, kickoutDo.OperatorCode)
	if err != nil {
		log.Error("[GroupController.KickOutMember]", "try to quit member", kickoutDo.MemberCode, "into group", kickoutDo.Code, ", appear an exception:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	go services.AddOperateRecord(group.Code, kickoutDo.OperatorCode, db.GroupMemberDataType, db.RemoveActionType,
		"user ", kickoutDo.OperatorCode, " member ", kickoutDo.MemberCode, " quit from group ", group.Code)
	backDo.Code = responses.OKCode
	backDo.Msg = "成功"
}

// ChangeMember 修改成员的群组信息 /group/customize
func (controller *GroupController) ChangeMember(ctx iris.Context) {
	backDo := responses.NewPackagedRespModel()
	defer func() {
		if err := recover(); err != nil {
			log.Error("[GroupController.ChangeMember]", "fatal when try to change member info in group, the exception is", err)
			backDo.Code = responses.ErrorCode
			backDo.Msg = "系统异常"
		}
		ctx.JSON(backDo)
	}()
	var customizeDo requests.GroupCustomizeReqModel
	if err := ctx.ReadJSON(&customizeDo); err != nil {
		log.Error("[GroupController.ChangeMember]", "read json appear an exception, err:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	group, err := controller.groupService.ChangeMemberNickname(customizeDo.Code, customizeDo.MemberCode, customizeDo.GroupAlias, customizeDo.GroupNickname, customizeDo.OperatorCode)
	if err != nil {
		log.Error("[GroupController.ChangeMember]", "try to change member", customizeDo.MemberCode, "info in group", customizeDo.Code, ", appear an exception:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	go services.AddOperateRecord(group.Code, customizeDo.OperatorCode, db.GroupMemberDataType, db.UpdateActionType,
		"user ", customizeDo.OperatorCode, " change member ", customizeDo.MemberCode, " group ", group.Code)
	backDo.Code = responses.OKCode
	backDo.Msg = "成功"
}

// Dismiss 解散群组 /group/dismiss
func (controller *GroupController) Dismiss(ctx iris.Context) {
	backDo := responses.NewPackagedRespModel()
	defer func() {
		if err := recover(); err != nil {
			log.Error("[GroupController.Dismiss]", "fatal when try to dismiss group, the exception is", err)
			backDo.Code = responses.ErrorCode
			backDo.Msg = "系统异常"
		}
		ctx.JSON(backDo)
	}()
	var dismissDo requests.GroupDismissReqModel
	if err := ctx.ReadJSON(&dismissDo); err != nil {
		log.Error("[GroupController.Dismiss]", "read json appear an exception, err:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	group, err := controller.groupService.EraseGroup(dismissDo.Code, dismissDo.OperatorCode)
	if err != nil {
		log.Error("[GroupController.Dismiss]", "try to dismiss group", dismissDo.Code, ", appear an exception:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	go services.AddOperateRecord(group.Code, dismissDo.OperatorCode, db.GroupDataType, db.RemoveActionType,
		"user ", dismissDo.OperatorCode, " dismiss group ", group.Code)
	backDo.Code = responses.OKCode
	backDo.Msg = "成功"
}
