package services

import (
	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"github.com/datsfilipe/pkg/repositories"
	"go.uber.org/zap"
)

type UserService struct {
	core.Service
	log *zap.Logger
	repo *repositories.UserRepository
}

func (us *UserService) CreateUser(data interface{}) (interface{}, error) {
	user := data.(*models.User)

	if user == nil {
		return nil, us.GenError(us.InvalidData, user)
	}

	duplicate, err := us.repo.FindByEmail(user.Email)
	if err != nil {
		us.log.Warn("Warn", zap.Any("Warn", err.Error()))
	}

	if duplicate != nil {
		return nil, us.GenError(us.DuplicateError, user)
	}

	dbRecord, err := us.repo.CreateUser(user)
	if err != nil {
		us.log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecord, nil
}

func (us *UserService) DeleteUser(data interface{}) (interface{}, error) {
	userID, ok := data.(string)

	if !ok || userID == "" {
		return nil, us.GenError(us.InvalidData, nil)
	}

	_, err := us.repo.DeleteUser(userID)
	if err != nil {
		us.log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return nil, nil
}

func (us *UserService) CreateProfile(data interface{}) (interface{}, error) {
	profile := data.(*models.Profile)

	if profile == nil {
		return nil, us.GenError(us.InvalidData, profile)
	}

	if duplicate, err := us.repo.FindByUsername(profile.Username); err != nil {
		us.log.Warn("Warn", zap.Any("Warn", err.Error()))
	} else if duplicate.(models.Profile).Username == profile.Username {
		return nil, us.GenError(us.DuplicateError, profile)
	}

	dbRecord, err := us.repo.CreateProfile(profile)
	if err != nil {
		us.log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecord, nil
}

func (us *UserService) FindByEmail(data interface{}) (interface{}, error) {
	email := data.(string)

	dbRecord, err := us.repo.FindByEmail(email)
	if err != nil {
		us.log.Warn("Warn", zap.Any("Warn", err.Error()))
	}

	return dbRecord, nil
}

func (us *UserService) FindByID(data interface{}) (interface{}, error) {
	id := data.(string)

	dbRecord, err := us.repo.FindByID(id)
	if err != nil {
		us.log.Warn("Warn", zap.Any("Warn", err.Error()))
	}

	return dbRecord, nil
}

func (us *UserService) GetProfiles(data interface{}) (interface{}, error) {
	userID, ok := data.(string)

	if !ok || userID == "" {
		return nil, us.GenError(us.InvalidData, nil)
	}

	dbRecord, err := us.repo.GetProfiles(userID)
	if err != nil {
		us.log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecord, nil
}

func (us *UserService) GetProfile(data interface{}) (interface{}, error) {
	profileID, ok := data.(string)

	if !ok || profileID == "" {
		return nil, us.GenError(us.InvalidData, nil)
	}

	dbRecord, err := us.repo.GetProfile(profileID)
	if err != nil {
		us.log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecord, nil
}

func (us *UserService) GetDefaultProfiles(data interface{}) (interface{}, error) {
	userIDs, ok := data.([]string)

	if !ok || len(userIDs) == 0 {
		return nil, us.GenError(us.InvalidData, nil)
	}

	dbRecord, err := us.repo.GetDefaultProfiles(userIDs)
	if err != nil {
		us.log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecord, nil
}

func NewUserService(log *zap.Logger, repo *repositories.UserRepository) *UserService {
	return &UserService{
		Service: *core.NewService(),
		log: log,
		repo: repo,
	}
}
