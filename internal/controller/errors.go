package controller

import "github.com/yakupovdev/FoodStore/internal/model"

type HTTPError struct {
	Status   int
	Response model.ErrorResponse
}

var (
	ErrInvalidJSON        = HTTPError{Status: 400, Response: model.ErrorResponse{Error: "invalid request body"}}
	ErrEmptyFields        = HTTPError{Status: 400, Response: model.ErrorResponse{Error: "email and password are required"}}
	ErrInvalidCredentials = HTTPError{Status: 401, Response: model.ErrorResponse{Error: "invalid email or password"}}
	ErrUserExists         = HTTPError{Status: 409, Response: model.ErrorResponse{Error: "user already exists"}}
	ErrInternal           = HTTPError{Status: 500, Response: model.ErrorResponse{Error: "internal server error"}}
)
