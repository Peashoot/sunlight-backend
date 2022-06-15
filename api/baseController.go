package api

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/peashoot/sunlight/entity/do/responses"
)

type BaseController interface {
	Route(party iris.Party)
}

func Init(app *iris.Application) {
	// 统一异常处理
	app.OnAnyErrorCode(FireErrorCode)
	// 给所有API留一个中间件的口子
	party := app.Party("/")
	// 注册所有控制器
	controllers := []BaseController{
		NewCategoryController(),
		NewFileController(),
		NewGroupController(),
		NewLabelController(),
		NewUserController(),
	}
	// 注册控制器的路由
	for _, controller := range controllers {
		controller.Route(party)
	}
}

// FireErrorCode 将所有异常的情况用JSON的格式返回
func FireErrorCode(ctx *context.Context) {
	backDo := responses.NewPackagedRespModel()
	backDo.Code = responses.RespCode(ctx.GetStatusCode())
	if ok, err := ctx.GetErrPublic(); ok {
		backDo.Msg = err.Error()
		ctx.JSON(backDo)
		return
	}
	backDo.Msg = context.StatusText(ctx.GetStatusCode())
	ctx.JSON(backDo)
}
