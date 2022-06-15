package db

type AuthTokenRecord struct {
	BaseModel
	UserID uint   `gorm:"comment:用户ID"`             // 用户ID
	Token  string `gorm:"size:512;comment:鉴权Token"` // 鉴权Token
}
