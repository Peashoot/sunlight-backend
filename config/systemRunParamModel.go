package config

import "time"

type SystemRunParamModel struct {
	ID        uint      `gorm:"primaryKey;comment:记录ID;autoIncrement"` // 记录id，自增
	Name      string    `gorm:"size:50;comment:配置名称;index"`            // 配置名称
	Kind      string    `gorm:"size:20;comment:配置类型"`                  // 配置类型
	Value     string    `gorm:"size:1024;comment:配置值"`                 // 配置值
	CreatedAt time.Time `gorm:"autoCreateTime;comment:创建时间"`           // 创建日期
	UpdatedAt time.Time `gorm:"autoUpdateTime;comment:更新时间"`           // 更新日期
}
