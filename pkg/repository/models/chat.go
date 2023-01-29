package models

type Chat struct {
	Id    int    `json:"id"`
	Name  string `json:"name" form:"name"  binding:"required"`
	Types string `json:"types"`
}

type ChatUsers struct {
	Id     int `json:"id"`
	ChatId int `json:"chat_id"`
	UserId int `json:"user_id"`
}
