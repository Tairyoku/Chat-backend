package models

type Status struct {
	Id           int    `json:"id" db:"id"`
	SenderId     int    `json:"sender_id" form:"sender_id" binding:"required"`
	RecipientId  int    `json:"recipient_id" form:"recipient_id" binding:"required"`
	Relationship string `json:"relationship" form:"relationship" binding:"required"`
}
