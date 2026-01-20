package model

type EmailRequest struct {
	Email string `json:"email"`
}

type CodeRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type ResetPasswordRequest struct {
	NewPassword string `json:"new_password"`
}

type EmailResponse struct {
	Message string `json:"message"`
}
