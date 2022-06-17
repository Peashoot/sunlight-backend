package db

// GroupMembershipModel 群组关系
type GroupMembershipModel struct {
	BaseModel
	MemberCode string `gorm:"size:50;comment:成员编号;index"` // 成员编号
	GroupCode  string `gorm:"size:50;comment:组编号;index"`  // 组编号
	InviteWay  int    `gorm:"comment:进群方式"`               // 进群方式
	Nickname   string `gorm:"size:128;comment:组昵称"`       // 组昵称
	GroupAlias string `gorm:"size:128;comment:群组别名"`      // 群组别名
}
