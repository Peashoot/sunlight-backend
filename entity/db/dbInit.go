package db

import "github.com/peashoot/sunlight/config"

func init() {
	config.MysqlDB.AutoMigrate(
		&AuthTokenRecord{},
		&CategoryModel{},
		&EnumDicModel{},
		&FileStorageRecord{},
		&GroupMembershipModel{},
		&RoleModel{},
		&UserGroupModel{},
		&UserLabelModel{},
		&UserModel{},
		&UserOperateRecord{},
		&config.SystemRunParamModel{},
	)
	config.Load()
}
