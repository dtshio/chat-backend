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
		core.Response(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err = json.Unmarshal(raw, &user)
	if err != nil {
		core.Response(w, http.StatusBadRequest, "Invalid user data")
		return
	}

	err = json.Unmarshal(raw, &profile)
	if err != nil {
		core.Response(w, http.StatusBadRequest, "Invalid profile data")
		return
	}

	userService := services.NewUserService()

	newUser, err := userService.CreateUser(user)
	if err != nil {
		core.Response(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	profile.UserID = newUser.(models.User).ID

	_, err = userService.CreateProfile(profile)
	if err != nil {
		core.Response(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	tokenPayload, err := json.Marshal(newUser.(models.User).ID)
	if err != nil {
		core.Response(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	token := auth.SignToken(string(tokenPayload))

	core.Response(w, http.StatusCreated, token)
}

func (uc *UserController) HandleSignIn(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}

	var raw json.RawMessage
	err := json.NewDecoder(r.Body).Decode(&raw)
	if err != nil {
		core.Response(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err = json.Unmarshal(raw, &user)
	if err != nil {
		core.Response(w, http.StatusBadRequest, "Invalid user data")
		return
	}

	userService := services.NewUserService()

	userRecord, _ := userService.FindByEmail(user.Email)
	if userRecord == nil {
		core.Response(w, http.StatusNotFound, "User not found")
		return
	}

	if userRecord.(models.User).Password != user.Password {
		core.Response(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	tokenPayload, err := json.Marshal(userRecord.(models.User).ID)
	if err != nil {
		core.Response(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	token := auth.SignToken(string(tokenPayload))
	core.Response(w, http.StatusOK, token)
}

func NewUserController() *UserController {
	userController := &UserController{}

	return &UserController{
		Controller: *core.NewController([]core.ControllerMethod{
			userController.HandleSignUp,
			userController.HandleSignIn,
		}),
	}
}
