package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/datsfilipe/pkg/application/auth"
	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"github.com/datsfilipe/pkg/services"
	"gorm.io/gorm"
)

type UserController struct {
	core.Controller
	db *gorm.DB
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

	newUser, err := userService.CreateUser(uc.db, user)
	if err != nil {
		core.Response(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	profile.UserID = newUser.(models.User).ID

	_, err = userService.CreateProfile(uc.db, profile)
	if err != nil {
		core.Response(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	stringfyProfile, _ := json.Marshal(profile)
	stringfyUser, _ := json.Marshal(newUser)
	token := auth.SignToken(fmt.Sprint(newUser.(models.User).ID))
	res := fmt.Sprintf(`{"token": "%s", "user": %s, "profile": %s}`, token, stringfyUser, stringfyProfile)
	core.Response(w, http.StatusCreated, res)
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

	userRecord, _ := userService.FindByEmail(uc.db, user.Email)
	if userRecord == nil {
		core.Response(w, http.StatusNotFound, "User not found")
		return
	}

	if userRecord.(models.User).Password != user.Password {
		core.Response(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	profile := &models.Profile{}
	profile.UserID = models.BigInt(userRecord.(models.User).ID)
	profileRecords, err := userService.GetProfiles(uc.db, profile.UserID)

	if profileRecords == nil {
		core.Response(w, http.StatusInternalServerError, "Could not find any profiles")
		return
	}

	stringfyProfiles, _ := json.Marshal(profileRecords.([]models.Profile))
	stringfyUser, _ := json.Marshal(userRecord)
	token := auth.SignToken(fmt.Sprint(userRecord.(models.User).ID))
	res := fmt.Sprintf(`{"token": "%s", "user": %s, "profiles": %s}`, token, stringfyUser, stringfyProfiles)
	core.Response(w, http.StatusOK, res)
}

func NewUserController(db *gorm.DB) *UserController {
	userController := &UserController{
		db: db,
	}

	return &UserController{
		Controller: *core.NewController([]core.ControllerMethod{
			userController.HandleSignUp,
			userController.HandleSignIn,
		}),
		db: db,
	}
}
