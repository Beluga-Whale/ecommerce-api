package services

import (
	"errors"

	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	"github.com/Beluga-Whale/ecommerce-api/internal/repositories"
	"github.com/Beluga-Whale/ecommerce-api/internal/utils"
)

type UserServiceInterface interface {
	Register(user *models.User)error
}

type UserService struct {
	userRepo repositories.UserRepositoryInterface
	hashPassword utils.ComparePasswordInterface
}

func NewUserService(userRepo repositories.UserRepositoryInterface , hashPassword utils.ComparePasswordInterface) *UserService {
	return &UserService{userRepo: userRepo, hashPassword: hashPassword}
}

func (s *UserService) Register(user *models.User)error {
	// NOTE - เช็คว่ามี email เป็นค่าว่างไหม
	if user.Email == "" || user.Password == "" {
		return errors.New("Email and Password cannot be empty")
	}

	// NOTE - เช็คว่ามี email ซ้ำไหม
	existingUser, err := s.userRepo.GetUserByEmail(user.Email)
	if err != nil {
		return errors.New("Error checking for existing user")
	}


	if existingUser != nil {
		return errors.New("Email already exists")
	}

	// NOTE - สร้าง user ใหม่
	
	return  s.userRepo.CreateUser(user)

}