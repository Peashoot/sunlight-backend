package db

// EnumDicModel 特殊枚举字典表
type EnumDicModel struct {
	Name     string `gorm:"size:50;comment:名称;index"` // 名称
	Value    string `gorm:"size:512;comment:字典值"`     // 字典值
	TypeName string `gorm:"size:50;comment:类型名称"`     // 类型名称
}
