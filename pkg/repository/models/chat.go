package models

type Chat struct {
	Id   int
	Name string
}

type ChatUsers struct {
	Id     int
	ChatId int
	UserId int
}
