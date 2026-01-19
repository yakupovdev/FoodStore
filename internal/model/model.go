package model

import "time"

type User struct {
	ID        int
	Login     string
	Password  string
	Type      string
	Created   time.Time
	LastEnter time.Time
	Balance   float64
}

type UserData struct {
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Type     string  `json:"type"`
	Balance  float64 `json:"balance"`
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GmailRequest struct {
	Email string `json:"email"`
}

type CodeRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type ResetPasswordRequest struct {
	NewPassword string `json:"new_password"`
}
