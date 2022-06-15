package requests

type GroupCreateReqModel struct {
	BaseReqModel
	Name         string   `json:"name"`       // 群组名称
	OperatorCode string   `json:"opUserCode"` // 操作员
	Members      []string `json:"member"`     // 成员编码
}
