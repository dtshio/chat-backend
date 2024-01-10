package services

import (
	"github.com/datsfilipe/pkg/application/database"
	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/repositories"
)

type UserService struct {
	core.Service
}

func (us *UserService) CreateUser(data interface{}) (interface{}, error) {
	db, err := database.Open()
	if err != nil {
		return nil, err
	}

	userRepo := repositories.NewUserRepository()

	user, err := userRepo.Create(db, data)
	return user, err
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
