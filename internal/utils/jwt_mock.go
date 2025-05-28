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

func (m *JwtMock) GenerateJWT(email string) (string, error) {
	args := m.Called(email)
	return args.String(0),args.Error(1)
}


func (m *JwtMock) ParseJWT(tokenString string) (string, error) {
	args := m.Called(tokenString)
	return args.String(0),args.Error(1)
}
