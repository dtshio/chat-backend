package controllers

import (
	"net/http"

	"github.com/datsfilipe/pkg/core"
)

type UserController struct {
	core.Controller
}

func (uc *UserController) HandleSignUp(w http.ResponseWriter, r *http.Request) {
}

func NewUserController() *UserController {
	userController := &UserController{}

	return &UserController{
		Controller: *core.NewController([]core.ControllerMethod{
			userController.HandleSignUp,
		}),
	}
}
