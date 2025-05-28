package utils

import (
	"github.com/stretchr/testify/mock"
)

type JwtMock struct {
	mock.Mock
}

func NewJwtMock() *JwtMock {
	return &JwtMock{}
}

func (m *JwtMock) GenerateJWT(email string, role string) (string, error) {
	args := m.Called(email,role)
	return args.String(0),args.Error(1)
}


func (m *JwtMock)  ParseJWT(tokenString string) (*JWTClaims, error) {
	args := m.Called(tokenString)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*JWTClaims), args.Error(1)
}
