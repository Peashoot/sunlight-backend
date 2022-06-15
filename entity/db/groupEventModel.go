package db

import "time"

// GroupEventModel 群组事件（纪念日、生日、某个重大事件、日记等）
type GroupEventModel struct {
	BaseModel
	Name      string    `gorm:"size:50;comment:名称"`  // 名称
	Abstract  string    `gorm:"size:512;comment:摘要"` // 摘要
	EventTime time.Time `gorm:"comment:发生时间"`        // 发生时间
	Detail    string    `gorm:"comment:事件详情"`        // 事件详情
	Status    uint8     `gorm:"comment:事件状态"`        // 事件状态
}
