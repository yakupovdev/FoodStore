package dto

type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	UserType string `json:"user_type"`
	Balance  int64  `json:"balance"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	UserType string `json:"user_type"`
}

type RegisterOutput struct {
	Message  string `json:"message"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	UserType string `json:"user_type"`
	Balance  int64  `json:"balance"`
}

type LoginOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshAccessTokenOutput struct {
	AccessToken string `json:"access_token"`
}
