package model

type User struct {
	UserID   string `gorm:"primaryKey;required;unique" json:"user_id,omitempty"`
	Password string `gorm:"required" json:"password,omitempty"`
	Nickname string `gorm:"size:30" json:"nickname,omitempty"`
	Comment  string `gorm:"size:100" json:"comment,omitempty"`
}
