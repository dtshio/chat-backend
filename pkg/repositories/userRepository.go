package repositories

import (
	"fmt"

	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	core.Repository
}

func (ur *UserRepository) CreateUser(db *gorm.DB, data interface {}) (interface {}, error) {
	user, ok := data.(*models.User)
    if !ok {
        return nil, fmt.Errorf("Data is not a valid user")
    }

	err := db.Table("users").Create(user).Error
	if err != nil {
		return nil, fmt.Errorf("Error creating user: %v", err)
	}

	return *user, err
}

func (ur *UserRepository) CreateProfile(db *gorm.DB, data interface {}) (interface {}, error) {
	profile, ok := data.(*models.Profile)
	if !ok {
		return nil, fmt.Errorf("Data is not a valid profile")
	}

	err := db.Table("user_profiles").Create(profile).Error
	if err != nil {
		return nil, fmt.Errorf("Error creating profile: %v", err)
	}

	return *profile, err
}

func (ur *UserRepository) FindByEmail(db *gorm.DB, data interface {}) (interface {}, error) {
	email := data.(string)
	var user models.User
	err := db.Table("users").Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, nil
	}

	return user, err
}

func (ur *UserRepository) FindByUsername(db *gorm.DB, data interface {}) (interface {}, error) {
	username := data.(string)
	var profile models.Profile
	err := db.Table("user_profiles").Where("username = ?", username).First(&profile).Error
	if err != nil {
		return nil, nil
	}

	return profile, err
}

func NewUserRepository() *UserRepository {
	userRepo := &UserRepository{}
	return &UserRepository{
		Repository: *core.NewRepository(
			&models.User{},
			[]core.RepositoryMethod{
				userRepo.CreateUser,
				userRepo.CreateProfile,
				userRepo.FindByEmail,
				userRepo.FindByUsername,
			},
		),
	}
}
