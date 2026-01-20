package usecase

import (
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

func (eu *EmailUsecase) SendCodeByEmail(emailTo string) error {
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

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"foodstorewwgo@gmail.com",
		[]string{emailTo},
		msg,
	)

	if err != nil {
		return err
	}

	userID, err := eu.repo.GetUserIDByEmail(emailTo)

	if err != nil {
		return err
	}

	err = eu.repo.SaveRecoveryCode(userID, emailTo, code)

	if err != nil {
		return err
	}

	return nil
}

func (eu *EmailUsecase) VerifyCode(email, code string) (string, error) {
	userID, err := eu.repo.GetUserIDByEmail(email)

	if err != nil {
		return "", ErrUserNotFound
	}

	isValidCode, err := eu.repo.VerifyRecoveryCode(email, code)
	if err != nil {
		return "", ErrVerificationFailed
	}

	if !isValidCode {
		return "", ErrCodeIsNotValid
	}

	token, err := security.GenerateToken(userID)

	if err != nil {
		return "", ErrTokenGeneration
	}

	_ = eu.repo.DeleteRecoveryCode(email)

	return token, nil
}
