package dto

import "time"

type DeleteUserInput struct {
	UserID int64 `json:"user_id"`
}

type CreateAdminInput struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	SecretKey string `json:"secret_key"`
}

type CreateAdminOutput struct {
	UserID  int64  `json:"user_id"`
	Message string `json:"message"`
}

type DeleteUserOutput struct {
	Message string `json:"message"`
}

type GetLogsHistoryOutput struct {
	LogID            int64     `json:"log_id"`
	ClientID         int64     `json:"client_id"`
	SellerID         int64     `json:"seller_id"`
	TotalAmount      int64     `json:"total_amount"`
	CommissionAmount int64     `json:"commission_amount"`
	CreatedAt        time.Time `json:"created_at"`
}

type GetAllUsersOutput struct {
	UserID        int64     `json:"user_id"`
	Email         string    `json:"email"`
	UserType      string    `json:"user_type"`
	Balance       int64     `json:"balance"`
	CreatedAt     time.Time `json:"created_at"`
	LastEnteredAt time.Time `json:"last_entered_at"`
}

type DeleteExpiredSubscriptionOutput struct {
	Message string `json:"message"`
}
