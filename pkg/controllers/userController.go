package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/datsfilipe/pkg/application/auth"
	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"github.com/datsfilipe/pkg/services"
)

type UserController struct {
	core.Controller
}

func (uc *UserController) HandleSignUp(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	profile := &models.Profile{}

	var raw json.RawMessage
	err := json.NewDecoder(r.Body).Decode(&raw)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request body"))
		return
	}

	err = json.Unmarshal(raw, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid user data"))
		return
	}

	err = json.Unmarshal(raw, &profile)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid profile data"))
		return
	}

	userService := services.NewUserService()

	newUser, err := userService.CreateUser(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	profile.UserID = newUser.(models.User).ID

	_, err = userService.CreateProfile(profile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	tokenPayload, err := json.Marshal(newUser.(models.User))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	token := auth.SignToken(string(tokenPayload))

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(token))
}

func NewUserController() *UserController {
	userController := &UserController{}

	return &UserController{
		Controller: *core.NewController([]core.ControllerMethod{
			userController.HandleSignUp,
		}),
	}
}
