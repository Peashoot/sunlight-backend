package responses

type PackagedRespModel struct {
	Code RespCode    `json:"code"` // 错误编码
	Msg  string      `json:"msg"`  // 错误说明
	Data interface{} `json:"data"` // 返回实体
}

type RespCode int

const (
	OKCode           RespCode = 200
	CreatedCode      RespCode = 201
	AcceptedCode     RespCode = 202
	BadRequestCode   RespCode = 400
	UnauthorizedCode RespCode = 401
	ForbiddenCode    RespCode = 403
	NotFoundCode     RespCode = 404
	NotAllowedCode   RespCode = 405
	TimeoutCode      RespCode = 408
	ErrorCode        RespCode = 500
	SystemBusyCode   RespCode = 502
)

func NewPackagedRespModel() *PackagedRespModel {
	model := &PackagedRespModel{}
	model.Code = -1
	model.Msg = "uninitialized"
	return model
}
