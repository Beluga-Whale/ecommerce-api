package services

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	"github.com/Beluga-Whale/ecommerce-api/internal/repositories"
	"github.com/Beluga-Whale/ecommerce-api/internal/utils"
)

type UserServiceInterface interface {
	Register(user *models.User)error
	Login(user *models.User) (string,error)
}

type UserService struct {
	userRepo repositories.UserRepositoryInterface
	hashPassword utils.ComparePasswordInterface
	jwtUtil utils.JwtInterface
}

func NewUserService(userRepo repositories.UserRepositoryInterface , hashPassword utils.ComparePasswordInterface,jwtUtil utils.JwtInterface) *UserService {
	return &UserService{userRepo: userRepo, hashPassword: hashPassword, jwtUtil: jwtUtil}
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

func (s *UserService) Login(user *models.User) (string,error) {
	// NOTE - เช็คว่ามี email เป็นค่าว่างไหม
	if user.Email == "" || user.Password == "" {
		return "",errors.New("Email and Password cannot be empty")
	}

	// NOTE - เช็คว่ามี email นี้ใน ฐานข้อมูลไหม
	dbUser, err := s.userRepo.GetUserByEmail(user.Email)

	if err != nil || dbUser == nil {
		return "", errors.New("User not found")
	}

	err = s.hashPassword.ComparePassword(dbUser.Password, user.Password)

	if err != nil {
		return "", errors.New("Invalid email or password")
	}

	userIDStr := strconv.FormatUint(uint64(dbUser.ID), 10)
	fmt.Print(userIDStr)
	token, err  := s.jwtUtil.GenerateJWT(dbUser.Email, string(dbUser.Role), userIDStr)

	if err != nil {
		return "", errors.New("Error generating JWT token")
	}


	return token,nil
}