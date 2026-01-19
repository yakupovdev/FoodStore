package model

type UserData struct {
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Type     string  `json:"type"`
	Balance  float64 `json:"balance"`
}
