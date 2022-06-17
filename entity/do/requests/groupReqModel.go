package requests

// GroupCreateReqModel 创建群组信息请求实体
type GroupCreateReqModel struct {
	BaseReqModel
	Name         string   `json:"name"`       // 群组名称
	Description  string   `json:"desc"`       // 群组描述
	OperatorCode string   `json:"opUserCode"` // 操作员
	Members      []string `json:"member"`     // 成员编码
}

// GroupChangeReqModel 改变群组信息请求实体
type GroupChangeReqModel struct {
	BaseReqModel
	Code         string `json:"code"`       // 群组编号
	Name         string `json:"name"`       // 群组名称
	Description  string `json:"desc"`       // 群组描述
	OwnerCode    string `json:"ownerCode"`  // 群主编号
	OperatorCode string `json:"opUserCode"` // 操作员编号
}

// GroupInviteMemberReqModel 邀请组员请求实体
type GroupInviteMemberReqModel struct {
	BaseReqModel
	Code         string   `json:"code"`       // 群组编号
	InviteWay    int      `json:"way"`        // 加群方式
	MemberCodes  []string `json:"members"`    // 组员编号数组
	OperatorCode string   `json:"opUserCode"` // 操作员编号
}

// GroupKickoutMemberReqModel 移除组员请求实体
type GroupKickoutMemberReqModel struct {
	BaseReqModel
	Code         string `json:"code"`       // 群组编号
	MemberCode   string `json:"member"`     // 组员编号
	OperatorCode string `json:"opUserCode"` // 操作员编号
}

// GroupCustomizeReqModel 自定义组员信息请求实体
type GroupCustomizeReqModel struct {
	BaseReqModel
	Code          string `json:"code"`       // 群组编号
	MemberCode    string `json:"member"`     // 组员编号
	GroupAlias    string `json:"alias"`      // 群组别名
	GroupNickname string `json:"nickname"`   // 群组昵称
	OperatorCode  string `json:"opUserCode"` // 操作员编号
}

// GroupDismissReqModel 解散群组请求实体
type GroupDismissReqModel struct {
	BaseReqModel
	Code         string `json:"code"`       // 群组编号
	OperatorCode string `json:"opUserCode"` // 操作员编号
}
