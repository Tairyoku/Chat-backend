package models

import "time"

type Message struct {
	Id     int       `json:"id" db:"id"`
	ChatId int       `json:"chat_id"`
	Author int       `json:"author"`
	Text   string    `json:"text" form:"text"  binding:"required"`
	SentAt time.Time `json:"sent_at"`
	//db:"sent_at" gorm:"->"
}
