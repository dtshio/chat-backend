package services

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/datsfilipe/pkg/application/database"
	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"github.com/datsfilipe/pkg/repositories"
)

func GenerateBigIntID() *big.Int {
	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(130), nil).Sub(max, big.NewInt(1))

	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		panic(err)
	}

	return n
}

type UserService struct {
	core.Service
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

	user.ID = GenerateBigIntID().Uint64()

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

	profile.ID = GenerateBigIntID().Uint64()

	return userRepository.CreateProfile(db, profile)
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
