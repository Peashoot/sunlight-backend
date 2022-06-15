package db

// UserGroupModel 用户组
type UserGroupModel struct {
	BaseModel
	Name string `gorm:"size:50;comment:名称"` // 名称
}
