package requests

type BaseReqModel struct {
	Timestamp uint64 `json:"ts"`  // 请求时间戳
	UID       string `json:"uid"` // 请求id
}
