package utils

import (
	"github.com/Beluga-Whale/ecommerce-api/internal/utils"
	"github.com/stretchr/testify/mock"
)

type JwtMock struct {
	mock.Mock
}

func NewJwtMock() *JwtMock {
	return &JwtMock{}
}

func (m *JwtMock) GenerateJWT(email string, role string, userID string) (string, error) {
	args := m.Called(email, role, userID)
	return args.String(0), args.Error(1)
}

func (m *JwtMock) ParseJWT(tokenString string) (*utils.JWTClaims, error) {
	args := m.Called(tokenString)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*utils.JWTClaims), args.Error(1)
}
