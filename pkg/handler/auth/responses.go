package auth

type TokenResponse struct {
	Token string `json:"token"`
}

type IdResponse struct {
	Id string `json:"id"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type UsernameInput struct {
	Username string `json:"username"`
}
