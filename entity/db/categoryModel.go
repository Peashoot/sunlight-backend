package db

// CategoryModel 类别（事件分类）
type CategoryModel struct {
	BaseModel
	Name        string `gorm:"size:50;comment:名称;index;"` // 名称
	Description string `gorm:"size:8192;comment:描述"`      // 描述
	State       uint   `gorm:"comment:状态"`                // 状态
	OwnerGroup  string `gorm:"comment:归属群组编码"`            // 归属群组编码
}
