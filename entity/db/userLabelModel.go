package db

// UserLabelModel 用户标签
type UserLabelModel struct {
	BaseModel
	Name   string `gorm:"size:50;comment:名称;index"` // 名称
	Status uint   `gorm:"comment:状态"`               // 标签状态
}

const (
	LabelBanedStatus     uint = 1 // 审核不通过
	LabelExaminingStatus uint = 5 // 正在审核中
	LabelAvaliableStatus uint = 9 // 审核通过
)
