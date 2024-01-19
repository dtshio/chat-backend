package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/datsfilipe/pkg/application/auth"
	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"github.com/datsfilipe/pkg/services"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserController struct {
	core.Controller
	db *gorm.DB
	log *zap.Logger
}


func (uc *UserController) HandleSignUp(w http.ResponseWriter, r *http.Request) {
	if uc.IsAllowedMethod(r, []string{"POST"}) == false {
		uc.Response(w, http.StatusMethodNotAllowed, nil)
		return
	}

	payload := uc.GetPayload(r)

	if payload == nil {
		uc.Response(w, http.StatusBadRequest, nil)
		return
	}

	user := &models.User{}
	user.Email = payload["email"].(string)
	user.Password = payload["password"].(string)

	profile := &models.Profile{}
	profile.Username = payload["username"].(string)

	if user.Email == "" || user.Password == "" || profile.Username == "" {
		uc.Response(w, http.StatusBadRequest, nil)
		return
	}

	service := services.NewUserService()

	userRecord, err := service.CreateUser(uc.db, uc.log, user)
	if err != nil {
		uc.Response(w, http.StatusInternalServerError, err)
		return
	}


	profile.UserID = userRecord.(models.User).ID

	profileRecord, err := service.CreateProfile(uc.db, uc.log, profile)
	if err != nil {
		uc.Response(w, http.StatusInternalServerError, err)
		return
	}

	userJson, _ := json.Marshal(userRecord)
	profileJson, _ := json.Marshal(profileRecord)

	userID := fmt.Sprint(userRecord.(models.User).ID)
	token := userID + "." + auth.SignToken(userID)

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

	service := services.NewUserService()

	userRecord, err := service.FindByEmail(uc.db, uc.log, user.Email)
	if userRecord == nil {
		uc.Response(w, http.StatusNotFound, err)
		return
	}

	if user.VerifyPassword(user.Password, userRecord.(models.User).Password) == false {
		uc.Response(w, http.StatusUnauthorized, nil)
		return
	}

	profile := &models.Profile{}
	profile.UserID = models.BigInt(userRecord.(models.User).ID)

	profileRecords, err := service.GetProfiles(uc.db, uc.log, profile.UserID)
	if err != nil {
		uc.Response(w, http.StatusInternalServerError, err)
		return
	}

	profilesJson, _ := json.Marshal(profileRecords.([]models.Profile))
	userJson, _ := json.Marshal(userRecord)

	userID := fmt.Sprint(userRecord.(models.User).ID)
	token := userID + "." + auth.SignToken(userID)

	res := fmt.Sprintf(`{"token": "%s", "user": %s, "profiles": %s}`, token, userJson, profilesJson)

	uc.Response(w, http.StatusOK, res)
}

func NewUserController(db *gorm.DB, log *zap.Logger) *UserController {
	controller := &UserController{
		db: db,
		log: log,
	}

	return &UserController{
		Controller: *core.NewController([]core.ControllerMethod{
			controller.HandleSignUp,
			controller.HandleSignIn,
		}),
		db: db,
		log: log,
	}
}
