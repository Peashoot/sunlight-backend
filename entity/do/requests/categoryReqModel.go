package requests

type CategoryCreateReqModel struct {
	BaseReqModel
	Name         string `json:"name"`       // 类别名称
	GroupCode    string `json:"groupCode"`  // 群组编码
	Description  string `json:"desc"`       // 类别描述
	OperatorCode string `json:"opUserCode"` // 操作员编码
}

type CategoryChangeNameReqModel struct {
	BaseReqModel
	CategoryCode string `json:"code"`       // 类别编码
	NewName      string `json:"name"`       // 类别名称
	OperatorCode string `json:"opUserCode"` // 操作员编码
}

type CategoryChangeDescReqModel struct {
	BaseReqModel
	CategoryCode string `json:"code"`       // 类别编码
	NewDesc      string `json:"desc"`       // 类别描述
	OperatorCode string `json:"opUserCode"` // 操作员编码
}

type CategoryRemoveReqModel struct {
	BaseReqModel
	CategoryCode string `json:"code"`       // 类别编码
	OperatorCode string `json:"opUserCode"` // 操作员编码
}
