package repositories

import (
	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	core.Repository
	db *gorm.DB
}

func (ur *UserRepository) CreateUser(data interface {}) (interface {}, error) {
	user, ok := data.(*models.User)

    if !ok {
        return nil, ur.GenError(ur.InvalidData, user)
    }

	err := user.BeforeCreateRecord()
	if err != nil {
		return nil, ur.GenError(ur.InvalidData, user)
	}

	err = ur.db.Table("users").Create(user).Error
	if err != nil {
		return nil, ur.GenError(ur.CreateError, user)
	}

	return *user, nil
}

func (ur *UserRepository) CreateProfile(data interface {}) (interface {}, error) {
	profile, ok := data.(*models.Profile)

	if !ok {
		return nil, ur.GenError(ur.InvalidData, profile)
	}

	err := profile.BeforeCreateRecord(profile.ID)
	if err != nil {
		return nil, ur.GenError(ur.InvalidData, profile)
	}

	err = ur.db.Table("user_profiles").Create(profile).Error
	if err != nil {
		return nil, ur.GenError(ur.CreateError, profile)
	}

	return *profile, nil
}

func (ur *UserRepository) FindByEmail(data interface {}) (interface {}, error) {
	email := data.(string)

	var user models.User

	err := ur.db.Table("users").Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) FindByUsername(data interface {}) (interface {}, error) {
	username := data.(string)

	var profile models.Profile

	err := ur.db.Table("user_profiles").Where("username = ?", username).First(&profile).Error
	if err != nil {
		return nil, err
	}

	return profile, nil 
}

func (ur *UserRepository) GetProfiles(data interface {}) (interface {}, error) {
	userID := data.(string)

	var profiles []models.Profile

	err := ur.db.Table("user_profiles").Where("user_id = ?", userID).Find(&profiles).Error

	if err != nil {
		return nil, nil
	}

	return profiles, err
}


func (ur *UserRepository) GetProfile(data interface {}) (interface {}, error) {
	profileID := data.(string)

	var profile []models.Profile

	err := ur.db.Table("user_profiles").Where("id = ?", profileID).Find(&profile).Error
	if err != nil {
		return nil, nil
	}

	return profile, err
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		Repository: *core.NewRepository(
			&models.User{},
			[]core.RepositoryMethod{
				(&UserRepository{}).CreateUser,
				(&UserRepository{}).CreateProfile,
				(&UserRepository{}).FindByEmail,
				(&UserRepository{}).FindByUsername,
				(&UserRepository{}).GetProfiles,
			},
		),
		db: db,
	}
}
