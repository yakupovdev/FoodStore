package model

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	UserType string `json:"user_type"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	UserType string `json:"user_type"`
	Balance  int64  `json:"balance"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterResponse struct {
	Message  string `json:"message"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	UserType string `json:"user_type"`
	Balance  int64  `json:"balance"`
}

type RefreshAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}
