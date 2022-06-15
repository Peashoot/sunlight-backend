package responses

type GetFileUploadAuthRespModel struct {
	Token      string `json:"token"`  // Token
	Expiration string `json:"expire"` // 有效期
}
