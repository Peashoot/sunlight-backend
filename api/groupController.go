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
	auth.AuthNeedParty.Post("/auth/group/create", controller.CreateNew)
}

// CreateNew 创建 /auth/group/create
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
	group, err := controller.groupService.CreateNewGroup(createDo.Name, createDo.OperatorCode, createDo.Members...)
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
}

// ChangeName 改名 /auth/group/changeName
func (controller *GroupController) ChangeName(ctx iris.Context) {

}
