package model

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	UserType string `json:"usertype"`
	Balance  int64  `json:"balance"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RegisterResponse struct {
	Message string `json:"message"`
}
