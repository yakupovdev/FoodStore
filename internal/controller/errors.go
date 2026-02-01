package controller

import "github.com/yakupovdev/FoodStore/internal/model"

type HTTPError struct {
	Status   int
	Response model.ErrorResponse
}

var (
	ErrInvalidJSON          = HTTPError{Status: 400, Response: model.ErrorResponse{Error: "invalid request body"}}
	ErrEmptyFields          = HTTPError{Status: 400, Response: model.ErrorResponse{Error: "some fields are missing"}}
	ErrInvalidCredentials   = HTTPError{Status: 401, Response: model.ErrorResponse{Error: "invalid email or password"}}
	ErrUserExists           = HTTPError{Status: 409, Response: model.ErrorResponse{Error: "user already exists"}}
	ErrInternal             = HTTPError{Status: 500, Response: model.ErrorResponse{Error: "internal server error"}}
	ErrUserNotFound         = HTTPError{Status: 404, Response: model.ErrorResponse{Error: "user not found"}}
	ErrVerifyCodeIsNotValid = HTTPError{Status: 400, Response: model.ErrorResponse{Error: "verify code is not valid"}}
	ErrInvalidToken         = HTTPError{Status: 401, Response: model.ErrorResponse{Error: "invalid token"}}
	ErrInvalidUserType      = HTTPError{Status: 401, Response: model.ErrorResponse{Error: "invalid user type"}}
)
