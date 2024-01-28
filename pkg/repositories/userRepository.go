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

func (ur *UserRepository) DeleteUser(data interface {}) (interface {}, error) {
	userID := data.(string)

	err := ur.db.Table("users").Delete(&models.User{}, userID).Error
	if err != nil {
		return nil, ur.GenError(ur.DeleteError, userID)
	}

	return nil, nil
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

func (ur *UserRepository) FindByID(data interface {}) (interface {}, error) {
	id := data.(string)

	var user models.User

	err := ur.db.Table("users").Where("id = ?", id).First(&user).Error
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

func (ur *UserRepository) GetDefaultProfiles(data interface{}) (interface{}, error) {
    userIDs := data.([]string)

    var users []models.User
    if err := ur.db.Where("id IN (?)", userIDs).Find(&users).Error; err != nil {
        return nil, err
    }

    var profiles []models.Profile
	for _, user := range users {
        var profile models.Profile
        if err := ur.db.Where("id = ?", user.DefaultProfileID).First(&profile).Error; err != nil {
            return nil, err
        }

        profiles = append(profiles, profile)
    }

    return profiles, nil
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
