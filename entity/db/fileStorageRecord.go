package db

type FileStorageRecord struct {
	BaseModel
	Name         string `gorm:"size:127;comment:名称;index"`   // 名称
	Subffix      string `gorm:"size:20;comment:后缀名"`         // 后缀名
	FileType     uint8  `gorm:"comment:文件类型"`                // 文件类型
	FileURL      string `gorm:"size:255;comment:文件链接;index"` // 文件链接
	RealURL      string `gorm:"size:1023;comment:文件直链"`      // 文件直链
	CategoryCode string `gorm:"size:50;comment:类别编码"`        // 列表编码
}
