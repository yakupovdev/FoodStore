package entity

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/yakupovdev/FoodStore/internal/domain"
)

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	UserType     string
	Name         string
	Balance      int64
	CreatedAt    time.Time
	LastEnter    time.Time
}

func NewUser(email, password, userType, name string, balance int64) (*User, error) {
	if email == "" {
		return nil, domain.ErrEmptyEmail
	}
	if password == "" {
		return nil, domain.ErrEmptyPassword
	}
	if userType != "client" && userType != "seller" {
		return nil, domain.ErrInvalidUserType
	}
	if name == "" {
		return nil, domain.ErrEmptyName
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	now := time.Now()
	return &User{
		Email:        email,
		PasswordHash: string(hash),
		UserType:     userType,
		Name:         name,
		Balance:      balance,
		CreatedAt:    now,
		LastEnter:    now,
	}, nil
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return string(hash), nil
}

func RestoreUser(id int64, email, passwordHash, userType, name string, balance int64, createdAt, lastEnter time.Time) *User {
	return &User{
		ID:           id,
		Email:        email,
		PasswordHash: passwordHash,
		UserType:     userType,
		Name:         name,
		Balance:      balance,
		CreatedAt:    createdAt,
		LastEnter:    lastEnter,
	}
}

func (u *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) == nil
}

func (u *User) IsClient() bool {
	return u.UserType == "client"
}

func (u *User) IsSeller() bool {
	return u.UserType == "seller"
}
