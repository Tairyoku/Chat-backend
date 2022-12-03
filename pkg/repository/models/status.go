package models

type Status struct {
	Id           int    `json:"id" db:"id"`
	SenderId     int    `json:"sender_id"`
	RecipientId  int    `json:"recipient_id"`
	Relationship string `json:"relationship"`
}
