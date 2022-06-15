package db

// GroupMembershipModel 群组关系
type GroupMembershipModel struct {
	BaseModel
	MemberCode string `gorm:"comment:成员编号;index"` // 成员编号
	GroupCode  string `gorm:"comment:组编号;index"`  // 组编号
	MemberName string `gorm:"comment:组昵称"`        // 组昵称
}
