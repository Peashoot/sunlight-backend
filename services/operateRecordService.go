package services

import (
	"fmt"

	"github.com/peashoot/sunlight/config"
	"github.com/peashoot/sunlight/entity/db"
)

// AddOperateRecord 添加操作记录
func AddOperateRecord(userCode, opUserCode string, dtType db.OpDataType, acType db.OpActionType, args ...interface{}) {
	record := db.UserOperateRecord{
		BaseModel: db.BaseModel{
			Code:      userCode,
			CreatedBy: opUserCode,
		},
		OperateUserCode: opUserCode,
		DataType:        dtType,
		ActionType:      acType,
		Detail:          fmt.Sprint(args...),
	}
	config.MysqlDB.Create(&record)
}
