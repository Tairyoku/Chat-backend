package models

type User struct {
	Id       int    `json:"id" db:"id"`
	Name     string `json:"name" form:"name" binding:"required"`
	Username string `json:"username" form:"username"  binding:"required"`
	Password string `json:"password" gorm:"column:password_hash" form:"password"  binding:"required"`
}

type UsersRelationship struct {
}
