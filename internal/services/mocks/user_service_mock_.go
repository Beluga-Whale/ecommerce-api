package services

import (
	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	"github.com/stretchr/testify/mock"
)

type UserServiceMock struct {
	mock.Mock
}

func NewUserServiceMock() *UserServiceMock {
	return &UserServiceMock{}
}

func (m *UserServiceMock) Register(user *models.User)error{
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserServiceMock) Login(user *models.User) (string,error){
	args := m.Called(user)
	return args.String(0), args.Error(1)
}