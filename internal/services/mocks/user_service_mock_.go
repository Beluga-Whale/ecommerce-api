package services

import (
	"github.com/Beluga-Whale/ecommerce-api/internal/dto"
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

func (m *UserServiceMock) GetProfile(userIDUint uint) (*models.User,error) {
	args := m.Called(userIDUint)
	
	if user,ok := args.Get(0).(*models.User); ok {
		return user,nil
	}

	return nil,args.Error(1)
}

func (m *UserServiceMock) UpdateProfile(userID uint, req dto.UserUpdateProfileDTO)  error {
	args := m.Called(userID,req)

	return args.Error(0)
}