package repositories

import (
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
        return nil, ur.GenError(ur.InvalidData, user)
    }

	err := user.BeforeCreateRecord()
	if err != nil {
		return nil, ur.GenError(ur.InvalidData, user)
	}

	err1 := db.Table("users").Create(user).Error
	if err1 != nil {
		return nil, ur.GenError(ur.CreatingError, user)
	}

	return *user, nil
}

func (ur *UserRepository) CreateProfile(db *gorm.DB, data interface {}) (interface {}, error) {
	profile, ok := data.(*models.Profile)
	if !ok {
		return nil, ur.GenError(ur.InvalidData, profile)
	}

	err := profile.BeforeCreateRecord()
	if err != nil {
		return nil, ur.GenError(ur.InvalidData, profile)
	}

	err1 := db.Table("user_profiles").Create(profile).Error
	if err1 != nil {
		return nil, ur.GenError(ur.CreatingError, profile)
	}

	return *profile, nil
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

func (ur *UserRepository) GetProfiles(db *gorm.DB, data interface {}) (interface {}, error) {
	var profiles []models.Profile
	err := db.Table("user_profiles").Where("user_id = ?", data).Find(&profiles).Error
	if err != nil {
		return nil, nil
	}

	return profiles, err
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
				userRepo.GetProfiles,
			},
		),
	}
}
