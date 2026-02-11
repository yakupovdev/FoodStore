package entity

type Seller struct {
	ID      int64
	Name    string
	Email   string
	Balance int64
	Rating  float32
}

type PurchaseSubscriptionParams struct {
	ID    int64
	Price int64
}

func RestoreSeller(id int64, name, email string, balance int64, rating float32) *Seller {
	return &Seller{
		ID:      id,
		Name:    name,
		Email:   email,
		Balance: balance,
		Rating:  rating,
	}
}
