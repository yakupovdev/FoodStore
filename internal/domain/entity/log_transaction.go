package entity

import "time"

type LogTransaction struct {
	ID               int64
	ClientID         int64
	SellerID         int64
	TotalAmount      int64
	CommissionAmount int64
	CreatedAt        time.Time
}
