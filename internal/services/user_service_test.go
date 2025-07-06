package services_test

import (
	"errors"
	"testing"
	"time"

	"github.com/Beluga-Whale/ecommerce-api/internal/dto"
	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	repositories "github.com/Beluga-Whale/ecommerce-api/internal/repositories/mocks"
	"github.com/Beluga-Whale/ecommerce-api/internal/services"
	utils "github.com/Beluga-Whale/ecommerce-api/internal/utils/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegister(t *testing.T){
	t.Run("Register Success",func(t *testing.T) {
		user := &models.User{
			Email: "test@gmail.com",
			Password: "password",
			FirstName: "TEST",
			LastName:"TEST",
			Phone: "0999999999",
			BirthDate:  time.Date(2000, time.March, 15, 0, 0, 0, 0, time.UTC),

		}

		userRepo := repositories.NewUserRepositoryMock()
		hashPassword := utils.NewComparePassMock()
		jwtUtil := utils.NewJwtMock()

		userRepo.On("GetUserByEmail",mock.Anything).Return(nil,nil)

		userRepo.On("CreateUser",user).Return(nil)
		
		userService := services.NewUserService(userRepo,hashPassword,jwtUtil)

		err := userService.Register(user)
		
		assert.NoError(t,err)

		// NOTE - เช็คว่ามีการ Call function ไหม
		userRepo.AssertExpectations(t)
		
	})

	t.Run("Email or Password is empty", func(t *testing.T) {
		user := &models.User{
			Email: "", 
			Password: "",
			FirstName: "Test",
			LastName: "Test",
			Phone: "0999999999",
			BirthDate: time.Date(2000, time.March, 15, 0, 0, 0, 0, time.UTC),
		}

		userRepo := repositories.NewUserRepositoryMock()
		hashPassword := utils.NewComparePassMock()
		jwtUtil := utils.NewJwtMock()

		userService := services.NewUserService(userRepo, hashPassword, jwtUtil)

		err := userService.Register(user)

		assert.EqualError(t, err, "Email and Password cannot be empty")
	})

	t.Run("FirstName or Password is LastName", func(t *testing.T) {
		user := &models.User{
			Email: "test@gmai.com", 
			Password: "password",
			FirstName: "",
			LastName: "",
			Phone: "0999999999",
			BirthDate: time.Date(2000, time.March, 15, 0, 0, 0, 0, time.UTC),
		}

		userRepo := repositories.NewUserRepositoryMock()
		hashPassword := utils.NewComparePassMock()
		jwtUtil := utils.NewJwtMock()

		userService := services.NewUserService(userRepo, hashPassword, jwtUtil)

		err := userService.Register(user)

		assert.EqualError(t, err, "FirstName and LastName cannot be empty")
	})

	t.Run("Phone", func(t *testing.T) {
		user := &models.User{
			Email: "test@gmai.com", 
			Password: "password",
			FirstName: "d",
			LastName: "d",
			Phone: "",
			BirthDate: time.Date(2000, time.March, 15, 0, 0, 0, 0, time.UTC),
		}

		userRepo := repositories.NewUserRepositoryMock()
		hashPassword := utils.NewComparePassMock()
		jwtUtil := utils.NewJwtMock()

		userService := services.NewUserService(userRepo, hashPassword, jwtUtil)

		err := userService.Register(user)

		assert.EqualError(t, err, "Phone cannot to be empty")
	})

	t.Run("BirthDate", func(t *testing.T) {
		user := &models.User{
			Email: "test@gmai.com", 
			Password: "password",
			FirstName: "test",
			LastName: "lastest",
			Phone: "00838383",
			BirthDate: time.Time{},
		}

		userRepo := repositories.NewUserRepositoryMock()
		hashPassword := utils.NewComparePassMock()
		jwtUtil := utils.NewJwtMock()

		userService := services.NewUserService(userRepo, hashPassword, jwtUtil)

		err := userService.Register(user)

		assert.EqualError(t, err, "BirthData cannot to by empty")
	})

	t.Run("Error to check exists email",func(t *testing.T) {
		user := &models.User{
			Email: "test@gmail.com",
			Password: "password",
			FirstName: "TEST",
			LastName:"TEST",
			Phone: "0999999999",
			BirthDate:  time.Date(2000, time.March, 15, 0, 0, 0, 0, time.UTC),

		}

		userRepo := repositories.NewUserRepositoryMock()
		hashPassword := utils.NewComparePassMock()
		jwtUtil := utils.NewJwtMock()

		userRepo.On("GetUserByEmail",mock.Anything).Return(nil,errors.New("Error checking for existing user"))
		
		userService := services.NewUserService(userRepo,hashPassword,jwtUtil)

		err := userService.Register(user)
		
		assert.EqualError(t,err,"Error checking for existing user")

		// NOTE - เช็คว่ามีการ Call function ไหม
		userRepo.AssertExpectations(t)
	})

	t.Run("Email is existing",func(t *testing.T) {
		user := &models.User{
			Email: "test@gmail.com",
			Password: "password",
			FirstName: "TEST",
			LastName:"TEST",
			Phone: "0999999999",
			BirthDate:  time.Date(2000, time.March, 15, 0, 0, 0, 0, time.UTC),

		}

		userRepo := repositories.NewUserRepositoryMock()
		hashPassword := utils.NewComparePassMock()
		jwtUtil := utils.NewJwtMock()

		userRepo.On("GetUserByEmail",mock.Anything).Return(user,nil)
		
		userService := services.NewUserService(userRepo,hashPassword,jwtUtil)

		err := userService.Register(user)
		
		assert.EqualError(t,err,"Email already exists")

		// NOTE - เช็คว่ามีการ Call function ไหม
		userRepo.AssertExpectations(t)
	})

}

func TestLogin(t *testing.T){
	t.Run("Login Success",func(t *testing.T) {
		user := &models.User{
			
			Email: "test@gmail.com",
			Password: "password",
		}

		userRepo := repositories.NewUserRepositoryMock()
		hashPassword := utils.NewComparePassMock()
		jwtUtil := utils.NewJwtMock()

		userRepo.On("GetUserByEmail",mock.Anything).Return(user,nil)

		hashPassword.On("ComparePassword",mock.Anything,mock.Anything).Return(nil)
		jwtUtil.On("GenerateJWT", mock.Anything, mock.Anything, mock.Anything).Return("mocked-token", nil)

		userService := services.NewUserService(userRepo,hashPassword,jwtUtil)

		token,_,err := userService.Login(user)
		
		assert.NotEmpty(t,token)
		assert.NoError(t,err)

		// NOTE - เช็คว่ามีการ Call function ไหม
		userRepo.AssertExpectations(t)
		
	})

	t.Run("Email or Password is required",func(t *testing.T) {
		user := &models.User{
			Email: "",
			Password: "",
		}

		userRepo := repositories.NewUserRepositoryMock()
		hashPassword := utils.NewComparePassMock()
		jwtUtil := utils.NewJwtMock()

		userService := services.NewUserService(userRepo,hashPassword,jwtUtil)

		_,_,err := userService.Login(user)
		
		assert.EqualError(t,err,"Email and Password cannot be empty")

	})

	t.Run("User not found",func(t *testing.T) {
		user := &models.User{
			Email: "test@gmail.com",
			Password: "password",
		}

		userRepo := repositories.NewUserRepositoryMock()
		hashPassword := utils.NewComparePassMock()
		jwtUtil := utils.NewJwtMock()

		userRepo.On("GetUserByEmail",mock.Anything).Return(nil,errors.New("User not found"))

		userService := services.NewUserService(userRepo,hashPassword,jwtUtil)

		_,_,err := userService.Login(user)
		
		assert.EqualError(t,err,"User not found")
	})
	
	t.Run("Invalid email or password",func(t *testing.T) {
		user := &models.User{
			
			Email: "test@gmail.com",
			Password: "password",
		}

		userRepo := repositories.NewUserRepositoryMock()
		hashPassword := utils.NewComparePassMock()
		jwtUtil := utils.NewJwtMock()

		userRepo.On("GetUserByEmail",mock.Anything).Return(user,nil)

		hashPassword.On("ComparePassword",mock.Anything,mock.Anything).Return(errors.New("Invalid email or password"))

		userService := services.NewUserService(userRepo,hashPassword,jwtUtil)

		_,_,err := userService.Login(user)
		
		assert.EqualError(t,err,"Invalid email or password")
		
	})

	t.Run("Login Success",func(t *testing.T) {
		user := &models.User{
			
			Email: "test@gmail.com",
			Password: "password",
		}

		userRepo := repositories.NewUserRepositoryMock()
		hashPassword := utils.NewComparePassMock()
		jwtUtil := utils.NewJwtMock()

		userRepo.On("GetUserByEmail",mock.Anything).Return(user,nil)

		hashPassword.On("ComparePassword",mock.Anything,mock.Anything).Return(nil)
		jwtUtil.On("GenerateJWT", mock.Anything, mock.Anything, mock.Anything).Return("mocked-token", errors.New("Error generating JWT token"))

		userService := services.NewUserService(userRepo,hashPassword,jwtUtil)

		_,_,err := userService.Login(user)
		
		assert.EqualError(t,err,"Error generating JWT token")
	})
}

func TestGetProfile(t *testing.T) {
	t.Run("GetProfile Success",func(t *testing.T) {
		usermock := &models.User{
			Email: "test@gmail.com",
		}

		userRepo := repositories.NewUserRepositoryMock()
		hashPassword := utils.NewComparePassMock()
		jwtUtil := utils.NewJwtMock()

		userRepo.On("GetProfileByUserId",mock.Anything).Return(usermock,nil)

		userService := services.NewUserService(userRepo,hashPassword,jwtUtil)

		user,err := userService.GetProfile(1)
		
		assert.NoError(t,err)
		assert.NotEmpty(t,user)

		// NOTE - เช็คว่ามีการ Call function ไหม
		userRepo.AssertExpectations(t)
	})

	t.Run("Error to get profile",func(t *testing.T) {
		userRepo := repositories.NewUserRepositoryMock()
		hashPassword := utils.NewComparePassMock()
		jwtUtil := utils.NewJwtMock()

		userRepo.On("GetProfileByUserId",mock.Anything).Return(nil,errors.New("Error got get profile"))

		userService := services.NewUserService(userRepo,hashPassword,jwtUtil)

		user,err := userService.GetProfile(1)
		
		assert.EqualError(t,err,"Error got get profile")
		assert.Empty(t,user)

		// NOTE - เช็คว่ามีการ Call function ไหม
		userRepo.AssertExpectations(t)
	})
}

func TestUpdateProfile(t *testing.T) {
	t.Run("UpdateProfile Success",func(t *testing.T) {
		fn:= "TEST"
		ln:= "Test"
		phone:= "0988333333"
		avatar:="TEST"

		req := dto.UserUpdateProfileDTO{
			FirstName: &fn,
			LastName: &ln,
			Phone: &phone,
			BirthDate: &time.Time{},
			Avatar: &avatar,
		}

		userRepo := repositories.NewUserRepositoryMock()
		hashPassword := utils.NewComparePassMock()
		jwtUtil := utils.NewJwtMock()

		userRepo.On("UpdateProfile",mock.Anything,mock.Anything).Return(nil)

		userService := services.NewUserService(userRepo,hashPassword,jwtUtil)

		err := userService.UpdateProfile(1,req)
		
		assert.NoError(t,err)

		// NOTE - เช็คว่ามีการ Call function ไหม
		userRepo.AssertExpectations(t)
	})

	t.Run("Error got get profile",func(t *testing.T) {
		fn:= "TEST"
		ln:= "Test"
		phone:= "0988333333"
		avatar:="TEST"

		req := dto.UserUpdateProfileDTO{
			FirstName: &fn,
			LastName: &ln,
			Phone: &phone,
			BirthDate: &time.Time{},
			Avatar: &avatar,
		}

		userRepo := repositories.NewUserRepositoryMock()
		hashPassword := utils.NewComparePassMock()
		jwtUtil := utils.NewJwtMock()

		userRepo.On("UpdateProfile",mock.Anything,mock.Anything).Return(errors.New("Error got get profile"))

		userService := services.NewUserService(userRepo,hashPassword,jwtUtil)

		err := userService.UpdateProfile(1,req)
		
		assert.EqualError(t,err,"Error got get profile")

		// NOTE - เช็คว่ามีการ Call function ไหม
		userRepo.AssertExpectations(t)
	})
}