package messages

import "cmd/pkg/repository/models"

type IdResponse struct {
	Id string `json:"id"`
}

type MessageResponse struct {
	Message models.Message `json:"message"`
}

type MessageListResponse struct {
	List []models.User `json:"list"`
}

type TextInput struct {
	Text string `json:"text"`
}
