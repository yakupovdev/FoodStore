package security

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/yakupovdev/FoodStore/internal/domain/entity"
)

var ErrUnexpectedSigningMethod = errors.New("unexpected signing method error")

type Claims struct {
	UserID   int64  `json:"user_id"`
	UserType string `json:"user_type"`
	jwt.RegisteredClaims
}

type JWTService struct{}

func NewJWTService() *JWTService {
	return &JWTService{}
}

var tokenTypeToSecret = map[entity.TokenType]string{
	entity.AccessTokenType:   "JWT_SECRET_ACCESS",
	entity.RefreshTokenType:  "JWT_SECRET_REFRESH",
	entity.RecoveryTokenType: "JWT_SECRET_RECOVERY",
}

var tokenTypeToTTL = map[entity.TokenType]time.Duration{
	entity.AccessTokenType:   1 * time.Hour,
	entity.RefreshTokenType:  7 * 24 * time.Hour,
	entity.RecoveryTokenType: 1 * time.Minute,
}

func (s *JWTService) GenerateToken(userID int64, userType string, tokenType entity.TokenType) (string, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file,ETOOO PIZDAAAA")
		return "", err
	}

	secretKey := tokenTypeToSecret[tokenType]
	jwtSecret := []byte(os.Getenv(secretKey))

	ttl := tokenTypeToTTL[tokenType]
	claims := Claims{
		UserID:   userID,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func (s *JWTService) ParseToken(tokenStr string, tokenType entity.TokenType) (*entity.TokenClaims, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	secretKey := tokenTypeToSecret[tokenType]
	jwtSecret := []byte(os.Getenv(secretKey))

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(
		tokenStr,
		claims,
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrUnexpectedSigningMethod
			}
			return jwtSecret, nil
		},
	)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return &entity.TokenClaims{
		UserID:   claims.UserID,
		UserType: claims.UserType,
	}, nil
}
