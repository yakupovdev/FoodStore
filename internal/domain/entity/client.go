package entity

type Client struct {
	ID       int64
	Name     string
	Email    string
	UserType string
	Balance  int64
	Rating   float64
}

func RestoreClient(id int64, name, email, userType string, balance int64, rating float64) *Client {
	return &Client{
		ID:       id,
		Name:     name,
		Email:    email,
		UserType: userType,
		Balance:  balance,
		Rating:   rating,
	}
}
