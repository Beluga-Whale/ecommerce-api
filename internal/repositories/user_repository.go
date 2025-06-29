package repositories

import (
	"github.com/Beluga-Whale/ecommerce-api/internal/dto"
	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetProfileByUserId(userIDUint uint) (*models.User, error)
	UpdateProfile(userID uint, req dto.UserUpdateProfileDTO) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository{
	return &UserRepository{db:db}
}

func (r *UserRepository) CreateUser(user *models.User)error {
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	return r.db.Create(user).Error
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No user found
		}
		return nil, err // Other error
	}
	return &user, nil // User found
}

func (r *UserRepository) GetProfileByUserId(userIDUint uint) (*models.User, error) {
	var user models.User

	err := r.db.Where("id = ?", userIDUint).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil 
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateProfile(userID uint, req dto.UserUpdateProfileDTO) error {
	updates := make(map[string]interface{})

	if req.FirstName != nil {
		updates["first_name"] = *req.FirstName
	}
	if req.LastName != nil {
		updates["last_name"] = *req.LastName
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.BirthDate != nil {
		updates["birth_date"] = *req.BirthDate
	}
	if req.Avatar != nil {
		updates["profile_image"] = *req.Avatar
	}

	if len(updates) == 0 {
		return nil 
	}

	return r.db.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error
}