package services

import (
	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"github.com/datsfilipe/pkg/repositories"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserService struct {
	core.Service
}

func (us *UserService) CreateUser(db *gorm.DB, log *zap.Logger, data interface{}) (interface{}, error) {
	user := data.(*models.User)

	if user == nil {
		return nil, us.GenError(us.InvalidData, user)
	}

	repo := repositories.NewUserRepository()

	duplicate, err := repo.FindByEmail(db, user.Email)
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
	}

	if duplicate != nil {
		return nil, us.GenError(us.DuplicateError, user)
	}

	dbRecord, err := repo.CreateUser(db, user)
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecord, nil
}

func (us *UserService) CreateProfile(db *gorm.DB, log *zap.Logger, data interface{}) (interface{}, error) {
	profile := data.(*models.Profile)

	if profile == nil {
		return nil, us.GenError(us.InvalidData, profile)
	}

	repo := repositories.NewUserRepository()

	duplicate, err := repo.FindByUsername(db, profile.Username)
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
	}

	if duplicate != nil {
		return nil, us.GenError(us.DuplicateError, profile)
	}

	dbRecord, err := repo.CreateProfile(db, profile)
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecord, nil
}

func (us *UserService) FindByEmail(db *gorm.DB, log *zap.Logger, data interface{}) (interface{}, error) {
	email := data.(string)

	repo := repositories.NewUserRepository()

	dbRecord, err := repo.FindByEmail(db, email)
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
	}

	return dbRecord, nil
}

func (us *UserService) GetProfiles(db *gorm.DB, log *zap.Logger, data interface{}) (interface{}, error) {
	repo := repositories.NewUserRepository()

	dbRecord, err := repo.GetProfiles(db, data)
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecord, nil
}

func NewUserService() *UserService {
	service := &UserService{}

	return &UserService{
		Service: *core.NewService(
			[]core.ServiceMethod{
				service.CreateUser,
				service.FindByEmail,
				service.CreateProfile,
				service.GetProfiles,
			},
		),
	}
}
