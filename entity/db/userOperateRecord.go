package db

// UserOperateRecord 用户操作日志
type UserOperateRecord struct {
	BaseModel
	OperateUserCode string       `gorm:"comment:操作员编号;index"` // 操作员编号
	DataType        OpDataType   `gorm:"comment:数据类型;index"`  // 数据类型
	ActionType      OpActionType `gorm:"comment:操作类型;index"`  // 操作类型
	Detail          string       `gorm:"comment:操作详情"`        // 操作详情
}

// OpDataType 数据类型
type OpDataType uint8

const (
	UserDataType OpDataType = iota
	LabelDataType
	CategoryDataType
	GroupDataType
	GroupMemberDataType
)

// OpActionType 操作类型
type OpActionType uint8

const (
	_                 OpActionType = iota
	InsertActionType               // 新增
	UpdateActionType               // 修改
	RemoveActionType               // 删除
	SelectActionType               // 查询
	CancelActionType               // 取消
	ConfirmActionType              // 确认
	LoadActionType                 // 加载
	CloseActionType                // 关闭
)
