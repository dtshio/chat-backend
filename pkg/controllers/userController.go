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
	if uc.IsAllowedMethod(r, []string{"POST"}) == false {
		uc.Response(w, http.StatusMethodNotAllowed, nil)
		return
	}

	user := &models.User{}
	profile := &models.Profile{}
	payload := uc.GetPayload(r)

	if payload == nil {
		uc.Response(w, http.StatusBadRequest, nil)
		return
	}

	user.Email = payload["email"].(string)
	user.Password = payload["password"].(string)
	profile.Username = payload["username"].(string)

	if user.Email == "" || user.Password == "" || profile.Username == "" {
		uc.Response(w, http.StatusBadRequest, nil)
		return
	}

	userService := services.NewUserService()
	newUser, err := userService.CreateUser(uc.db, user)

	if err != nil {
		uc.Response(w, http.StatusInternalServerError, nil)
		return
	}

	profile.UserID = newUser.(models.User).ID

	_, err = userService.CreateProfile(uc.db, profile)
	if err != nil {
		uc.Response(w, http.StatusInternalServerError, nil)
		return
	}

	profileJson, _ := json.Marshal(profile)
	userJson, _ := json.Marshal(newUser)

	id := fmt.Sprint(newUser.(models.User).ID)
	token := id + "." + auth.SignToken(id)

	res := fmt.Sprintf(`{"token": "%s", "user": %s, "profile": %s}`, token, userJson, profileJson)

	uc.Response(w, http.StatusCreated, res)
}

func (uc *UserController) HandleSignIn(w http.ResponseWriter, r *http.Request) {
	if uc.IsAllowedMethod(r, []string{"POST"}) == false {
		uc.Response(w, http.StatusMethodNotAllowed, nil)
		return
	}

	user := &models.User{}
	payload := uc.GetPayload(r)

	if payload == nil {
		uc.Response(w, http.StatusBadRequest, nil)
		return
	}

	user.Email = payload["email"].(string)
	user.Password = payload["password"].(string)

	if user.Email == "" || user.Password == "" {
		uc.Response(w, http.StatusBadRequest, nil)
		return
	}

	userService := services.NewUserService()

	userRecord, _ := userService.FindByEmail(uc.db, user.Email)
	if userRecord == nil {
		uc.Response(w, http.StatusNotFound, nil)
		return
	}

	if user.VerifyPassword(user.Password, userRecord.(models.User).Password) == false {
		uc.Response(w, http.StatusUnauthorized, nil)
		return
	}

	profile := &models.Profile{}
	profile.UserID = models.BigInt(userRecord.(models.User).ID)
	profileRecords, _ := userService.GetProfiles(uc.db, profile.UserID)

	if profileRecords == nil {
		uc.Response(w, http.StatusInternalServerError, nil)
		return
	}

	profilesJson, _ := json.Marshal(profileRecords.([]models.Profile))
	userJson, _ := json.Marshal(userRecord)

	id := fmt.Sprint(userRecord.(models.User).ID)
	token := id + "." + auth.SignToken(id)

	res := fmt.Sprintf(`{"token": "%s", "user": %s, "profiles": %s}`, token, userJson, profilesJson)

	uc.Response(w, http.StatusOK, res)
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
