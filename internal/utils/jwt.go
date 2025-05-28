package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtInterface interface {
	GenerateJWT(email string) (string, error)
	ParseJWT(tokenString string) (string, error)
}

type JWTClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func NewJwt() *JWTClaims{
	return &JWTClaims{}
}

func (c *JWTClaims) GenerateJWT(email string) (string, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET"))

	claims :=JWTClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)

}

func (c *JWTClaims) ParseJWT(tokenString string) (string, error) {
	// NOTE - นำ token มาเช็คว่าเป็นอันเดียวกันไหม
	token, err := jwt.ParseWithClaims(tokenString,&JWTClaims{}, func(token *jwt.Token)  (interface{},error){
		return []byte(os.Getenv("JWT_SECRET")),nil
	})

	if err !=nil || !token.Valid {
		return "" , err
	}

	// NOTE - ดึงข้อมูลจาก claim
	claims, ok := token.Claims.(*JWTClaims); 

	if !ok {
		return "",errors.New("Invalid token claims")
	}

	return claims.Email,nil

}