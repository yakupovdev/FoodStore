package handler

import (
	"github.com/yakupovdev/FoodStore/internal/delivery/http/dto"
)

type HTTPError struct {
	Status   int
	Response dto.ErrorOutput
}

var (
	ErrInvalidJSON          = HTTPError{Status: 400, Response: dto.ErrorOutput{Error: "invalid request body"}}
	ErrEmptyFields          = HTTPError{Status: 400, Response: dto.ErrorOutput{Error: "some fields are missing"}}
	ErrInvalidCredentials   = HTTPError{Status: 401, Response: dto.ErrorOutput{Error: "invalid email or password"}}
	ErrUserExists           = HTTPError{Status: 409, Response: dto.ErrorOutput{Error: "user already exists"}}
	ErrInternal             = HTTPError{Status: 500, Response: dto.ErrorOutput{Error: "internal server error"}}
	ErrUserNotFound         = HTTPError{Status: 404, Response: dto.ErrorOutput{Error: "user not found"}}
	ErrVerifyCodeIsNotValid = HTTPError{Status: 400, Response: dto.ErrorOutput{Error: "verify code is not valid"}}
	ErrInvalidToken         = HTTPError{Status: 401, Response: dto.ErrorOutput{Error: "invalid token"}}
	ErrInvalidUserType      = HTTPError{Status: 401, Response: dto.ErrorOutput{Error: "invalid user type"}}
	ErrInvalidData          = HTTPError{Status: 400, Response: dto.ErrorOutput{Error: "invalid data"}}
	ErrInvalidID            = HTTPError{Status: 400, Response: dto.ErrorOutput{Error: "invalid ID"}}
	ErrNoProducts           = HTTPError{Status: 404, Response: dto.ErrorOutput{Error: "no products found"}}
	ErrOfferNotFound        = HTTPError{Status: 404, Response: dto.ErrorOutput{Error: "offer not found"}}
	ErrNotEnoughQuantity    = HTTPError{Status: 400, Response: dto.ErrorOutput{Error: "not enough quantity available"}}
	ErrNotEnoughBalance     = HTTPError{Status: 400, Response: dto.ErrorOutput{Error: "not enough balance to complete the subscription"}}
	ErrSubscriptionNotFound = HTTPError{Status: 400, Response: dto.ErrorOutput{Error: "subscription failed, try again later,check your balance or contact support"}}
)
