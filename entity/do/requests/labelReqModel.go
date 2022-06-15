package requests

// UserLabelCreateReqModel 标签创建请求实体
type UserLabelCreateReqModel struct {
	BaseReqModel
	Info     string `json:"info"`     // 标签内容
	UserCode string `json:"userCode"` // 用户编号
}

// UserLabelAbandonReqModel 标签删除请求实体
type UserLabelAbandonReqModel struct {
	BaseReqModel
	LabelCode string `json:"code"`      // 标签编号
	UserCode  string `json:"userCode "` // 用户编号
}

// UserLabelExamineReqModel 标签审核请求实体
type UserLabelExamineReqModel struct {
	BaseReqModel
	LabelCode string `json:"code"`     // 标签编号
	UserCode  string `json:"userCode"` // 用户编号
	Status    uint   `json:"status"`   // 标签状态
}
