package security

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

type TokenType string

const (
	AccessToken   TokenType = "JWT_SECRET_ACCESS"
	RefreshToken  TokenType = "JWT_SECRET_REFRESH"
	RecoveryToken TokenType = "JWT_SECRET_RECOVERY"

	AccessTokenLiving   time.Duration = 1 * time.Hour
	RefreshTokenLiving  time.Duration = 7 * time.Hour
	RecoveryTokenLiving time.Duration = 1 * time.Hour
)

var mapTokenLiving = map[TokenType]time.Duration{
	AccessToken:   AccessTokenLiving,
	RefreshToken:  RefreshTokenLiving,
	RecoveryToken: RecoveryTokenLiving,
}

func loadEnv() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	return nil
}

func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	hashHex := hex.EncodeToString(hash[:])
	return hashHex
}

func GenerateToken(userID int64, typeToken TokenType) (string, error) {
	if err := loadEnv(); err != nil {
		return "", err
	}
	jwtSecret := []byte(os.Getenv(string(typeToken)))

	date := jwt.NewNumericDate(time.Now().Add(mapTokenLiving[typeToken]))
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: date,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseToken(tokenStr string, typeToken TokenType) (*Claims, error) {
	if err := loadEnv(); err != nil {
		return nil, err
	}
	jwtSecret := []byte(os.Getenv(string(typeToken)))
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

	return claims, nil
}

func GenerateAccessCodeByEmail() string {
	number := rand.Intn(900000) + 100000
	str := strconv.Itoa(number)
	return str
}
