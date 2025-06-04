package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtInterface interface {
	GenerateJWT(email string,role string, userId string) (string, error)
	ParseJWT(tokenString string) (*JWTClaims, error)
}

type JWTClaims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	UserID string `json:"userID"`
	jwt.RegisteredClaims
}

func NewJwt() *JWTClaims{
	return &JWTClaims{}
}

func (c *JWTClaims) GenerateJWT(email string,role string, userId string) (string, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET"))

	claims :=JWTClaims{
		Email: email,
		Role:  role,
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)

}

func (c *JWTClaims) ParseJWT(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New("Invalid token claims")
	}

	return claims, nil
}
