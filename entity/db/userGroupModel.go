package db

// UserGroupModel 用户组
type UserGroupModel struct {
	BaseModel
	Name        string `gorm:"size:50;comment:名称;index"`   // 名称
	Description string `gorm:"size:512;comment:描述"`        // 描述
	OwnerCode   string `gorm:"size:50;comment:群主编号;index"` // 群主编号
}
