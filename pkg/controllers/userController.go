package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/datsfilipe/pkg/application/auth"
	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"github.com/datsfilipe/pkg/services"
)

type UserController struct {
	core.Controller
	service *services.UserService
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

	userRecord, err := uc.service.CreateUser(user)
	if err != nil {
		uc.Response(w, http.StatusInternalServerError, err)
		return
	}

	profile.UserID = userRecord.(models.User).ID
	profile.ID = userRecord.(models.User).DefaultProfileID

	profileRecord, err := uc.service.CreateProfile(profile)
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

	userRecord, err := uc.service.FindByEmail(user.Email)
	if userRecord == nil {
		uc.Response(w, http.StatusNotFound, err)
		return
	}

	if user.VerifyPassword(user.Password, userRecord.(models.User).Password) == false {
		uc.Response(w, http.StatusUnauthorized, nil)
		return
	}

	profile := &models.Profile{}
	profile.UserID = userRecord.(models.User).ID
	userID := fmt.Sprint(userRecord.(models.User).ID)

	profileRecords, err := uc.service.GetProfiles(userID)
	if err != nil {
		uc.Response(w, http.StatusInternalServerError, err)
		return
	}

	profilesJson, _ := json.Marshal(profileRecords.([]models.Profile))
	userJson, _ := json.Marshal(userRecord)

	token := userID + "." + auth.SignToken(userID)

	res := fmt.Sprintf(`{"token": "%s", "user": %s, "profiles": %s}`, token, userJson, profilesJson)

	uc.Response(w, http.StatusOK, res)
}

func (uc *UserController) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	if uc.IsAllowedMethod(r, []string{"DELETE"}) == false {
		uc.Response(w, http.StatusMethodNotAllowed, nil)
		return
	}

	userID := uc.GetPayload(r)["id"].(string)

	if userID == "" {
		uc.Response(w, http.StatusBadRequest, nil)
		return
	}

	if _, err := uc.service.DeleteUser(userID); err != nil {
		uc.Response(w, http.StatusInternalServerError, err)
		return
	}

	uc.Response(w, http.StatusOK, nil)
}

func (uc *UserController) HandleCreateProfile(w http.ResponseWriter, r *http.Request) {
	if uc.IsAllowedMethod(r, []string{"POST"}) == false {
		uc.Response(w, http.StatusMethodNotAllowed, nil)
		return
	}

	token := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]
	userIDStr := strings.Split(token, ".")[0]

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		uc.Response(w, http.StatusBadRequest, nil)
		return
	}

	payload := uc.GetPayload(r)

	profile := &models.Profile{}
	profile.UserID = models.BigInt(userID)
	profile.Username = payload["username"].(string)

	dbRecord, err := uc.service.CreateProfile(profile)
	if err != nil {
		uc.Response(w, http.StatusInternalServerError, err)
		return
	}

	res, err := json.Marshal(dbRecord)
	if err != nil {
		uc.Response(w, http.StatusInternalServerError, nil)
		return
	}

	uc.Response(w, http.StatusCreated, res)
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{
		service: service,
	}
}
