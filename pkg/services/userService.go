package services

import (
	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"github.com/datsfilipe/pkg/repositories"
	"gorm.io/gorm"
)

type UserService struct {
	core.Service
}

func (us *UserService) CreateUser(db *gorm.DB, data interface{}) (interface{}, error) {
	user := data.(*models.User)
	if user == nil {
		return nil, us.GenError(us.InvalidData, user)
	}

	userRepo := repositories.NewUserRepository()
	userRecord, _ := userRepo.FindByEmail(db, user.Email)

	if userRecord != nil {
		return nil, us.GenError(us.DuplicateError, user)
	}

	user.ID = models.BigInt(core.GenerateID())

	return userRepo.CreateUser(db, user)
}

func (us *UserService) CreateProfile(db *gorm.DB, data interface{}) (interface{}, error) {
	profile := data.(*models.Profile)
	if profile == nil {
		return nil, us.GenError(us.InvalidData, profile)
	}

	userRepository := repositories.NewUserRepository()
	profileRecord, _ := userRepository.FindByUsername(db, profile.Username)

	if profileRecord != nil {
		return nil, us.GenError(us.DuplicateError, profile)
	}

	profile.ID = models.BigInt(core.GenerateID())

	return userRepository.CreateProfile(db, profile)
}

func (us *UserService) FindByEmail(db *gorm.DB, data interface{}) (interface{}, error) {
	email := data.(string)

	userRepository := repositories.NewUserRepository()

	return userRepository.FindByEmail(db, email)
}

func (us *UserService) GetProfiles(db *gorm.DB, data interface{}) (interface{}, error) {
	userRepository := repositories.NewUserRepository()

	return userRepository.GetProfiles(db, data)
}

func NewUserService() *UserService {
	userService := &UserService{}

	return &UserService{
		Service: *core.NewService(
			[]core.ServiceMethod{
				userService.CreateUser,
				userService.FindByEmail,
				userService.CreateProfile,
				userService.GetProfiles,
			},
		),
	}
}
