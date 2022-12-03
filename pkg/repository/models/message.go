package models

import "time"

type Message struct {
	Id     int
	ChatId int
	UserId int
	Body   string
	SentAt time.Time
}
