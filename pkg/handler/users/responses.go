package users

import "cmd/pkg/repository/models"

type IdResponse struct {
	Id string `json:"id"`
}

type UserResponse struct {
	User models.User `json:"user"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ListResponse struct {
	List []models.User `json:"list"`
}

type StatusesListResponse struct {
	Friends     []models.User `json:"friends"`
	Blacklist   []models.User `json:"blacklist"`
	OnBlacklist []models.User `json:"onBlacklist"`
	Invites     []models.User `json:"invites"`
	Requires    []models.User `json:"requires"`
}
