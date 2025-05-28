package services_test

import (
	"errors"
	"testing"

	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	"github.com/Beluga-Whale/ecommerce-api/internal/repositories"
	"github.com/Beluga-Whale/ecommerce-api/internal/services"
	"github.com/Beluga-Whale/ecommerce-api/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegister(t *testing.T) {
	t.Run("Create User Successfully", func(t *testing.T) {
		user := &models.User{
			Email: "test@gmail.com",
			Password: "password123",
		}
		
		userRepo := repositories.NewUserRepositoryMock()
		hashPassword := utils.NewComparePassMock()

		userRepo.On("GetUserByEmail", mock.Anything).Return(nil,nil)
		userRepo.On("CreateUser", mock.Anything).Return(nil)

		userService := services.NewUserService(userRepo,hashPassword)

		err := userService.Register(user)

		assert.NoError(t,err)

		userRepo.AssertExpectations(t)
	})

	t.Run("Email Empty", func(t *testing.T) {
		user := &models.User{
			Email: "",
			Password: "password123",
		}

		userRepo := repositories.NewUserRepositoryMock()
		hashPassword := utils.NewComparePassMock()
		userService := services.NewUserService(userRepo,hashPassword)

		err := userService.Register(user)

		assert.EqualError(t,err,"Email and Password cannot be empty")
	})
	t.Run("Email Empty", func(t *testing.T) {
		user := &models.User{
			Email: "test@gmail.com",
			Password: "",
		}

		userRepo := repositories.NewUserRepositoryMock()
		hashPassword := utils.NewComparePassMock()
		userService := services.NewUserService(userRepo,hashPassword)

		err := userService.Register(user)

		assert.EqualError(t,err,"Email and Password cannot be empty")
	})
	t.Run("Fail to check email", func(t *testing.T) {
		user := &models.User{
			Email: "test@mail.com",
			Password: "password123",
		}

		userRepo := repositories.NewUserRepositoryMock()
		userRepo.On("GetUserByEmail", user.Email).Return(user, errors.New("Error checking for existing user"))

		hashPassword := utils.NewComparePassMock()
		userService := services.NewUserService(userRepo,hashPassword)
		err := userService.Register(user)
		assert.EqualError(t, err,"Error checking for existing user")
	})

	t.Run("Email already exists", func(t *testing.T) {
		user := &models.User{
			Email: "test@gmail.com",
			Password: "password123",
		}
		
		userRepo := repositories.NewUserRepositoryMock()

		userRepo.On("GetUserByEmail", "test@gmail.com").Return(user,nil)

		hashPassword := utils.NewComparePassMock()
		userService := services.NewUserService(userRepo,hashPassword)

		err := userService.Register(user)

		assert.EqualError(t,err,"Email already exists")

	})
}