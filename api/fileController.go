package api

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"github.com/peashoot/sunlight/auth"
	"github.com/peashoot/sunlight/entity/do/responses"
	"github.com/peashoot/sunlight/log"
	"github.com/peashoot/sunlight/utils"
)

type FileController struct {
}

func NewFileController() *FileController {
	return &FileController{}
}

func (controller *FileController) Route(app iris.Party) {
	auth.AuthNeedParty.Get("/auth/file/getoken", controller.GetUploadToken)
}

// GetUploadToken 获取文件上传Token /auth/file/getoken
func (controller *FileController) GetUploadToken(ctx iris.Context) {
	baseResp := responses.NewPackagedRespModel()
	defer func() {
		if err := recover(); err != nil {
			log.Error("[FileController.GetUploadToken]", "fatal when try to get token, the exception is", err)
			baseResp.Code = responses.ErrorCode
			baseResp.Msg = "系统异常"
		}
		ctx.JSON(baseResp)
	}()
	cred, err := utils.GetAliOss().GetUploadCredentials()
	if err != nil {
		log.Error("[FileController.GetUploadToken]", "try to get token, appear an exception:", err)
		baseResp.Code = responses.ErrorCode
		baseResp.Msg = "系统异常"
	}
	inner := responses.GetFileUploadAuthRespModel{}
	inner.Token = cred.SecurityToken
	inner.Expiration = cred.Expiration
	claims := jwt.Get(ctx).(*auth.AuthClaims)
	log.Info("[FileController.GetUploadToken]", "user", claims.Code, "get file upload token", cred.SecurityToken, ", the expiration is", cred.Expiration)
	baseResp.Data = inner
	baseResp.Code = responses.OKCode
	baseResp.Msg = "成功"
}
