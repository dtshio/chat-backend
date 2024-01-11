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

func GenerateUniqueBigIntID(existingIDs map[string]bool) *big.Int {
	maxInt64 := new(big.Int)
	maxInt64.Exp(big.NewInt(2), big.NewInt(63), nil).Sub(maxInt64, big.NewInt(1))

	for {
		n, err := rand.Int(rand.Reader, maxInt64)
		if err != nil {
			continue
		}

		idStr := n.String()
		if !existingIDs[idStr] {
			return n
		}
	}
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

	existingIDs := make(map[string]bool)
	user.ID = GenerateUniqueBigIntID(existingIDs).Uint64()

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

	existingIDs := make(map[string]bool)
	profile.ID = GenerateUniqueBigIntID(existingIDs).Uint64()

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
