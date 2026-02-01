package usecase

import (
	"errors"
	"net/smtp"

	"github.com/yakupovdev/FoodStore/internal/repository"
	"github.com/yakupovdev/FoodStore/internal/security"
)

type EmailUsecase struct {
	repo *repository.Postgres
}

func NewEmailUsecase(repo *repository.Postgres) (*EmailUsecase, error) {
	if repo == nil {
		return nil, ErrDatabaseConnection
	}

	return &EmailUsecase{
		repo: repo,
	}, nil
}

func (eu *EmailUsecase) SendCodeByEmail(emailTo, userType string) error {
	userID, err := eu.repo.GetUserIDByEmailAndType(emailTo, userType)

	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return repository.ErrUserNotFound
		} else {
			return ErrInternalServer
		}
	}

	code := security.GenerateAccessCodeByEmail()

	auth := smtp.PlainAuth(
		"",
		"foodstorewwgo@gmail.com",
		"dkeywmbvieuiuazj",
		"smtp.gmail.com",
	)

	msg := []byte(
		"From: FoodStore <foodstorewwgo@gmail.com>\r\n" +
			"To:" + emailTo + "\r\n" +
			"Subject: Test Email from FoodStore\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-UserType: text/plain; charset=\"UTF-8\"\r\n" +
			"\r\n" +
			"Here your code: " + code + "\n",
	)

	err = smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"foodstorewwgo@gmail.com",
		[]string{emailTo},
		msg,
	)

	if err != nil {
		return ErrSMTPFailed
	}

	err = eu.repo.SaveRecoveryCode(userID, emailTo, userType, code)

	if err != nil {
		return ErrInternalServer
	}

	return nil
}

func (eu *EmailUsecase) VerifyCode(email, userType, code string) (string, error) {
	userID, err := eu.repo.GetUserIDByEmailAndType(email, userType)

	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", repository.ErrUserNotFound
		} else {
			return "", ErrInternalServer
		}
	}

	isValidCode, err := eu.repo.VerifyRecoveryCode(email, userType, code)
	if err != nil {
		return "", ErrVerificationFailed
	}

	if !isValidCode {
		return "", ErrCodeIsNotValid
	}

	recoveryToken, err := security.GenerateToken(userID, userType, security.RecoveryToken)

	if err != nil {
		return "", ErrTokenGeneration
	}

	_ = eu.repo.DeleteRecoveryCode(email, userType)

	return recoveryToken, nil
}
