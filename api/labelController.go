package api

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/peashoot/sunlight/auth"
	"github.com/peashoot/sunlight/entity/db"
	"github.com/peashoot/sunlight/entity/do/requests"
	"github.com/peashoot/sunlight/entity/do/responses"
	"github.com/peashoot/sunlight/log"
	"github.com/peashoot/sunlight/services"
)

// LabelController 用户标签控制器
type LabelController struct {
	labelService *services.LabelService
}

func NewLabelController() *LabelController {
	return &LabelController{
		labelService: services.NewLabelService(),
	}
}

func (controller *LabelController) Route(app iris.Party) {
	auth.AuthNeedParty.Post("/label/create", controller.AddNew)
	auth.AuthNeedParty.Post("/label/abandon", controller.Abandon)
	auth.AuthNeedParty.Post("/label/examine", controller.Examine)
}

// AddNew 添加新标签 /label/create
func (controller *LabelController) AddNew(ctx iris.Context) {
	backDo := responses.NewPackagedRespModel()
	defer func() {
		if err := recover(); err != nil {
			log.Error("[LabelController.AddNew]", "fatal when try to create, the exception is", err)
			backDo.Code = responses.ErrorCode
			backDo.Msg = "系统异常"
		}
		ctx.JSON(backDo)
	}()
	var createDo requests.UserLabelCreateReqModel
	if err := ctx.ReadJSON(&createDo); err != nil {
		log.Error("[LabelController.AddNew]", "read json appear an exception, err:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	label, err := controller.labelService.AddOrGetExists(createDo.Info, createDo.UserCode)
	if err != nil {
		log.Error("[LabelController.AddNew]", "try to add new label", createDo.Info, ", appear an exception:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	go services.AddOperateRecord(label.Code, createDo.UserCode, db.LabelDataType, db.InsertActionType,
		"user", createDo.UserCode, "create a label", label.Code, "info:", label.Name)
	backDo.Code = responses.OKCode
	backDo.Msg = "成功"
	backDo.Data = &responses.UserLabelCreateRespModel{
		Code: label.Code,
	}
}

// Abandon 丢弃新标签 /label/abandon
func (controller *LabelController) Abandon(ctx iris.Context) {
	backDo := responses.NewPackagedRespModel()
	defer func() {
		if err := recover(); err != nil {
			log.Error("[LabelController.Abandon]", "fatal when try to remove, the exception is", err)
			backDo.Code = responses.ErrorCode
			backDo.Msg = "系统异常"
		}
		ctx.JSON(backDo)
	}()
	var removeDo requests.UserLabelAbandonReqModel
	if err := ctx.ReadJSON(&removeDo); err != nil {
		log.Error("[LabelController.Abandon]", "read json appear an exception, err:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	label, err := controller.labelService.Remove(removeDo.LabelCode, removeDo.UserCode)
	if err != nil {
		log.Error("[LabelController.Abandon]", "try to remove label", removeDo.LabelCode, ", appear an exception:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	go services.AddOperateRecord(label.Code, removeDo.UserCode, db.LabelDataType, db.RemoveActionType,
		"user", removeDo.UserCode, "remove a label", label.Code, "info:", label.Name)
	backDo.Code = responses.OKCode
	backDo.Msg = "成功"
}

// Examine 审核标签 /label/examine
func (controller *LabelController) Examine(ctx iris.Context) {
	backDo := responses.NewPackagedRespModel()
	defer func() {
		if err := recover(); err != nil {
			log.Error("[LabelController.Examine]", "fatal when try to examine, the exception is", err)
			backDo.Code = responses.ErrorCode
			backDo.Msg = "系统异常"
		}
		ctx.JSON(backDo)
	}()
	var examineDo requests.UserLabelExamineReqModel
	if err := ctx.ReadJSON(&examineDo); err != nil {
		log.Error("[LabelController.Examine]", "read json appear an exception, err:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	label, err := controller.labelService.Examine(examineDo.LabelCode, examineDo.UserCode, examineDo.Status)
	if err != nil {
		log.Error("[LabelController.Examine]", "try to examine label", examineDo.LabelCode, ", appear an exception:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	go services.AddOperateRecord(label.Code, examineDo.UserCode, db.LabelDataType, db.UpdateActionType,
		"user", examineDo.UserCode, "examine a label", label.Name, "status:", strconv.Itoa(int(examineDo.Status)))
	backDo.Code = responses.OKCode
	backDo.Msg = "成功"
}
