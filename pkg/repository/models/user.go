package models

type User struct {
	Id       int    `json:"id" db:"id"`
	Username string `json:"username" form:"username"  binding:"required"`
	Password string `json:"password" gorm:"column:password_hash" form:"password"  binding:"required"`
	Icon     string `json:"icon" binding:"required"`
}
