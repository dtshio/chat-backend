package services

import (
	"fmt"
	"time"

	"github.com/datsfilipe/pkg/application/database"
	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"github.com/datsfilipe/pkg/repositories"
	"github.com/disgoorg/snowflake/v2"
)

type UserService struct {
	core.Service
}

func GenerateID() uint64 {
	return uint64(snowflake.New(time.Now()))
}

func (us *UserService) CreateUser(data interface{}) (interface{}, error) {
	user := data.(*models.User)
	if user == nil {
		return nil, fmt.Errorf("Data is not a valid user")
	}

	db, err := database.Open()
	if err != nil {
		return nil, err
	}

	userRepo := repositories.NewUserRepository()
	userRecord, _ := userRepo.FindByEmail(db, user.Email)

	if userRecord != nil {
		return nil, fmt.Errorf("User already exists")
	}

	user.ID = GenerateID()

	return userRepo.CreateUser(db, user)
}

func (us *UserService) CreateProfile(data interface{}) (interface{}, error) {
	profile := data.(*models.Profile)
	if profile == nil {
		return nil, fmt.Errorf("Data is not a valid profile")
	}

	db, err := database.Open()
	if err != nil {
		return nil, err
	}

	userRepository := repositories.NewUserRepository()
	profileRecord, _ := userRepository.FindByUsername(db, profile.Username)

	if profileRecord != nil {
		return nil, fmt.Errorf("Profile already exists")
	}

	profile.ID = GenerateID()

	return userRepository.CreateProfile(db, profile)
}

func (us *UserService) FindByEmail(data interface{}) (interface{}, error) {
	email := data.(string)

	db, err := database.Open()
	if err != nil {
		return nil, err
	}

	userRepository := repositories.NewUserRepository()

	return userRepository.FindByEmail(db, email)
}

func NewUserService() *UserService {
	userService := &UserService{}

	return &UserService{
		Service: *core.NewService(
			[]core.ServiceMethod{
				userService.CreateUser,
			},
		),
	}
}
