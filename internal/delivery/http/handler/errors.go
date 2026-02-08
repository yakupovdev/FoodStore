package handler

import "github.com/yakupovdev/FoodStore/internal/usecase/dto"

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
)
