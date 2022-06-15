package db

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 数据库实体基础字段
type BaseModel struct {
	ID        uint           `gorm:"primaryKey;comment:记录ID;autoIncrement"` // 记录id，自增
	Code      string         `gorm:"size:50;comment:编号;index"`              // 编号，用于记录部分初始化信息字段
	CreatedBy string         `gorm:"comment:创建人编号"`                         // 创建人编号
	UpdatedBy string         `gorm:"comment:更新人编号"`                         // 更新人编号
	DeletedBy string         `gorm:"comment:删除人编号"`                         // 删除人编号
	CreatedAt time.Time      `gorm:"autoCreateTime;comment:创建时间"`           // 创建日期
	UpdatedAt time.Time      `gorm:"comment:更新时间"`                          // 更新日期
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间"`                    // 删除日期
}
