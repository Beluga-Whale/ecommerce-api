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
		jwtUtils := utils.NewJwtMock()

		userRepo.On("GetUserByEmail", mock.Anything).Return(nil,nil)
		userRepo.On("CreateUser", mock.Anything).Return(nil)

		userService := services.NewUserService(userRepo,hashPassword,jwtUtils)

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
		jwtUtils := utils.NewJwtMock()
		userService := services.NewUserService(userRepo,hashPassword,jwtUtils)

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
		jwtUtils := utils.NewJwtMock()
		userService := services.NewUserService(userRepo,hashPassword,jwtUtils)

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
		jwtUtils := utils.NewJwtMock()
		userService := services.NewUserService(userRepo,hashPassword,jwtUtils)
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
		jwtUtils := utils.NewJwtMock()
		userService := services.NewUserService(userRepo,hashPassword,jwtUtils)

		err := userService.Register(user)

		assert.EqualError(t,err,"Email already exists")

	})
}

func TestLogin(t *testing.T) {
	t.Run("Login Success",func(t *testing.T) {

		dbUser := &models.User{
			Email: "test@gmail.com",
			Password: "hash_password",
		}

		inputUser := &models.User{
			Email: "test@gmail.com",
			Password: "password",
		}

		userRepo := repositories.NewUserRepositoryMock()
		hashPassword := utils.NewComparePassMock()
		jwtUtils := utils.NewJwtMock()

		userRepo.On("GetUserByEmail",inputUser.Email).Return(dbUser,nil)

		hashPassword.On("ComparePassword", dbUser.Password, inputUser.Password).Return(nil)

		jwtUtils.On("GenerateJWT", dbUser.Email).Return("jwt_token", nil)

		userService := services.NewUserService(userRepo,hashPassword,jwtUtils)

		token, err := userService.Login(inputUser)

		assert.NoError(t, err)
		assert.Equal(t, "jwt_token", token)
		userRepo.AssertExpectations(t)
		
	})
	t.Run("User not found",func(t *testing.T) {

		inputUser := &models.User{
			Email: "",
			Password: "password",
		}

		userRepo := repositories.NewUserRepositoryMock()
		hashPassword := utils.NewComparePassMock()
		jwtUtils := utils.NewJwtMock()

		userService := services.NewUserService(userRepo,hashPassword,jwtUtils)

		_, err := userService.Login(inputUser)

		assert.EqualError(t, err,"Email and Password cannot be empty")
		
	})

	t.Run("User not found",func(t *testing.T) {

		dbUser := &models.User{
			Email: "test@gmail.com",
			Password: "hash_password",
		}

		inputUser := &models.User{
			Email: "test@gmail.com",
			Password: "password",
		}

		userRepo := repositories.NewUserRepositoryMock()
		hashPassword := utils.NewComparePassMock()
		jwtUtils := utils.NewJwtMock()

		userRepo.On("GetUserByEmail",inputUser.Email).Return(dbUser,errors.New("User not found"))

		userService := services.NewUserService(userRepo,hashPassword,jwtUtils)

		_, err := userService.Login(inputUser)

		assert.EqualError(t, err,"User not found")

		userRepo.AssertExpectations(t)
		
	})

	t.Run("Invalid email or password",func(t *testing.T) {

		dbUser := &models.User{
			Email: "test@gmail.com",
			Password: "hash_password",
		}

		inputUser := &models.User{
			Email: "test@gmail.com",
			Password: "password",
		}

		userRepo := repositories.NewUserRepositoryMock()
		hashPassword := utils.NewComparePassMock()
		jwtUtils := utils.NewJwtMock()

		userRepo.On("GetUserByEmail",inputUser.Email).Return(dbUser,nil)

		hashPassword.On("ComparePassword", dbUser.Password, inputUser.Password).Return(errors.New("Invalid email or password"))

		userService := services.NewUserService(userRepo,hashPassword,jwtUtils)

		_, err := userService.Login(inputUser)

		assert.EqualError(t, err,"Invalid email or password")

		userRepo.AssertExpectations(t)
		
	})

	t.Run("Login Success",func(t *testing.T) {

		dbUser := &models.User{
			Email: "test@gmail.com",
			Password: "hash_password",
		}

		inputUser := &models.User{
			Email: "test@gmail.com",
			Password: "password",
		}

		userRepo := repositories.NewUserRepositoryMock()
		hashPassword := utils.NewComparePassMock()
		jwtUtils := utils.NewJwtMock()

		userRepo.On("GetUserByEmail",inputUser.Email).Return(dbUser,nil)

		hashPassword.On("ComparePassword", dbUser.Password, inputUser.Password).Return(nil)

		jwtUtils.On("GenerateJWT", dbUser.Email).Return("", errors.New("Error generating JWT token"))

		userService := services.NewUserService(userRepo,hashPassword,jwtUtils)

		token, err := userService.Login(inputUser)

		assert.EqualError(t, err, "Error generating JWT token")
		assert.Equal(t, "", token)
		userRepo.AssertExpectations(t)
		
	})
}