package db

// RoleModel 用户角色
type RoleModel struct {
	BaseModel
	Name       string `gorm:"size:50;comment:角色名称;index"` // 角色名称
	PowerValue string `gorm:"size:256;comment:权限二进制序列"`   // 权限值二进制序列（每一位上的1代表有权限，0代表没有权限）
}
