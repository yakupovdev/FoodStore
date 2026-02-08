package repository

import "context"

type TransactionRepository interface {
	ExecuteOrderTransaction(ctx context.Context, clientID int64, totalAmount int64) error

	ExecuteSellerTransaction(ctx context.Context, sellerID int64, totalAmount int64) error
}
