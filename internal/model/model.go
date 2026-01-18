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
	Login    string  `json:"login"`
	Password string  `json:"password"`
	Type     string  `json:"type"`
	Balance  float64 `json:"balance"`
}

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
