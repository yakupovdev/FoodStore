package dto

type SendCodeInput struct {
	Email    string `json:"email"`
	UserType string `json:"user_type"`
}

type VerifyCodeInput struct {
	Email    string `json:"email"`
	Code     string `json:"code"`
	UserType string `json:"user_type"`
}

type ResetPasswordInput struct {
	UserID      int64
	NewPassword string `json:"new_password"`
}

type SendCodeOutput struct {
	Message string `json:"message"`
}

type VerifyCodeOutput struct {
	RecoveryToken string `json:"recovery_token"`
	Message       string `json:"message"`
}

type ResetPasswordOutput struct {
	Message string `json:"message"`
}

type ErrorOutput struct {
	Error string `json:"error"`
}
