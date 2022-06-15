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

type CategoryController struct {
	categoryService *services.CategoryService
}

func NewCategoryController() *CategoryController {
	return &CategoryController{
		categoryService: &services.CategoryService{},
	}
}

func (controller *CategoryController) Route(app iris.Party) {
	auth.AuthNeedParty.Post("/auth/category/create", controller.CreateNew)
	auth.AuthNeedParty.Post("/auth/category/changeName", controller.ChangeName)
	auth.AuthNeedParty.Post("/auth/category/changeDesc", controller.ChangeDesc)
	auth.AuthNeedParty.Post("/auth/category/remove", controller.Remove)
}

// CreateNew 创建新类别 /auth/category/create
func (controller *CategoryController) CreateNew(ctx iris.Context) {
	backDo := responses.NewPackagedRespModel()
	defer func() {
		if err := recover(); err != nil {
			log.Error("[CategoryController.CreateNew]", "fatal when try to create, the exception is", err)
			backDo.Code = responses.ErrorCode
			backDo.Msg = "系统异常"
		}
		ctx.JSON(backDo)
	}()
	var createDo requests.CategoryCreateReqModel
	if err := ctx.ReadJSON(&createDo); err != nil {
		log.Error("[CategoryController.CreateNew]", "read json appear an exception, err:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	category, err := controller.categoryService.Create(createDo.Name, createDo.GroupCode, createDo.OperatorCode, createDo.Description)
	if err != nil {
		log.Error("[CategoryController.CreateNew]", "try to add new category", createDo.Name, ", appear an exception:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	go services.AddOperateRecord(category.Code, createDo.OperatorCode, db.CategoryDataType, db.InsertActionType,
		"user", createDo.OperatorCode, "create a category", category.Code, "info:", category.Name)
	backDo.Code = responses.OKCode
	backDo.Msg = "成功"
}

// ChangeName 修改种类名称 /auth/category/changeName
func (controller *CategoryController) ChangeName(ctx iris.Context) {
	backDo := responses.NewPackagedRespModel()
	defer func() {
		if err := recover(); err != nil {
			log.Error("[CategoryController.ChangeName]", "fatal when try to change name, the exception is", err)
			backDo.Code = responses.ErrorCode
			backDo.Msg = "系统异常"
		}
		ctx.JSON(backDo)
	}()
	var changeDo requests.CategoryChangeNameReqModel
	if err := ctx.ReadJSON(&changeDo); err != nil {
		log.Error("[CategoryController.ChangeName]", "read json appear an exception, err:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	category, err := controller.categoryService.ChangeName(changeDo.NewName, changeDo.CategoryCode, changeDo.OperatorCode)
	if err != nil {
		log.Error("[CategoryController.ChangeName]", "try to change category", changeDo.CategoryCode, "'s name", changeDo.NewName, ", appear an exception:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	go services.AddOperateRecord(category.Code, changeDo.OperatorCode, db.CategoryDataType, db.UpdateActionType,
		"user", changeDo.OperatorCode, "change category", category.Code, "'s name:", category.Name)
	backDo.Code = responses.OKCode
	backDo.Msg = "成功"
}

// ChangeDesc 修改描述 /auth/category/changeDesc
func (controller *CategoryController) ChangeDesc(ctx iris.Context) {
	backDo := responses.NewPackagedRespModel()
	defer func() {
		if err := recover(); err != nil {
			log.Error("[CategoryController.ChangeDesc]", "fatal when try to change desc, the exception is", err)
			backDo.Code = responses.ErrorCode
			backDo.Msg = "系统异常"
		}
		ctx.JSON(backDo)
	}()
	var changeDo requests.CategoryChangeDescReqModel
	if err := ctx.ReadJSON(&changeDo); err != nil {
		log.Error("[CategoryController.ChangeDesc]", "read json appear an exception, err:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	category, err := controller.categoryService.ChangeDesc(changeDo.NewDesc, changeDo.CategoryCode, changeDo.OperatorCode)
	if err != nil {
		log.Error("[CategoryController.ChangeDesc]", "try to change category", changeDo.CategoryCode, "'s desc", changeDo.NewDesc, ", appear an exception:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	go services.AddOperateRecord(category.Code, changeDo.OperatorCode, db.CategoryDataType, db.UpdateActionType,
		"user", changeDo.OperatorCode, "change category", category.Code, "'s desc:", category.Description)
	backDo.Code = responses.OKCode
	backDo.Msg = "成功"
}

// Remove 移除 /auth/category/remove
func (controller *CategoryController) Remove(ctx iris.Context) {
	backDo := responses.NewPackagedRespModel()
	defer func() {
		if err := recover(); err != nil {
			log.Error("[CategoryController.Remove]", "fatal when try to remove category, the exception is", err)
			backDo.Code = responses.ErrorCode
			backDo.Msg = "系统异常"
		}
		ctx.JSON(backDo)
	}()
	var removeDo requests.CategoryRemoveReqModel
	if err := ctx.ReadJSON(&removeDo); err != nil {
		log.Error("[CategoryController.Remove]", "read json appear an exception, err:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	category, err := controller.categoryService.Remove(removeDo.CategoryCode, removeDo.OperatorCode)
	if err != nil {
		log.Error("[CategoryController.Remove]", "try to remove category", removeDo.CategoryCode, ", appear an exception:", err)
		backDo.Code = responses.ErrorCode
		backDo.Msg = "系统异常"
		return
	}
	go services.AddOperateRecord(category.Code, removeDo.OperatorCode, db.CategoryDataType, db.RemoveActionType,
		"user", removeDo.OperatorCode, "remove category", category.Code)
	backDo.Code = responses.OKCode
	backDo.Msg = "成功"
}
