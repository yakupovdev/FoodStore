package repository

import "context"

type AdminRepository interface {
	UpdateBalance(ctx context.Context, userID int64, newBalance int64) error

	GetExpiringSubscriptions(ctx context.Context) ([]int64, error)

	CancelSubscription(ctx context.Context, sellerID int64) error

	SetPriorityToFalse(ctx context.Context, sellerID int64) error
}
