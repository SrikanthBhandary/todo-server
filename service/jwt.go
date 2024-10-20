package service

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTValidator interface {
	ValidateToken(tokenString string) (int, error)
	GenerateToken(userID int) (string, error)
}
type JWTService struct {
	secretKey string
}

func NewJWTService(secret string) JWTValidator {
	return &JWTService{secretKey: secret}
}

// GenerateToken generates a new JWT token for a user
func (j *JWTService) GenerateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// ValidateToken validates the JWT token and returns the user ID if valid
func (j *JWTService) ValidateToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure that the token's signing method is what you expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("error validating the token")
		}
		return []byte(j.secretKey), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return int(claims["sub"].(float64)), nil
	}
	return 0, err
}
