package model

type VerifyEmailRequest struct {
	Email string `json:"email"`
}

type VerifyCodeRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type ResetUserPasswordRequest struct {
	NewPassword string `json:"new_password"`
}

type VerifyEmailResponse struct {
	Message string `json:"message"`
}

type VerifyCodeResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

type ResetUserPasswordResponse struct {
	Message string `json:"message"`
}
