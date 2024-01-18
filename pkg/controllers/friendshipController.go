package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"github.com/datsfilipe/pkg/services"
	"gorm.io/gorm"
)

type FriendshipController struct {
	core.Controller
	db *gorm.DB
}

func (fc *FriendshipController) HandleNewFriendship(w http.ResponseWriter, r *http.Request) {
	if fc.IsAllowedMethod(r, []string{"POST"}) == false {
		fc.Response(w, http.StatusMethodNotAllowed, nil)
		return
	}

	if fc.IsAuthorized(r) == false {
		fc.Response(w, http.StatusUnauthorized, nil)
		return
	}

	payload := fc.GetPayload(r)

	initiatorID := payload["initiator_id"].(string)
	friendID := payload["friend_id"].(string)

	if initiatorID == "" || friendID == "" {
		fc.Response(w, http.StatusBadRequest, nil)
		return
	}

	friendshipService := services.NewFriendshipService()
	newFriendship, err := friendshipService.CreateFriendship(fc.db, payload)

	if err != nil {
		fc.Response(w, http.StatusInternalServerError, nil)
		return
	}

	res := fmt.Sprint("{\"id\": \"", newFriendship.(models.Friendship).ID, ", \"channel_id\": \"", newFriendship.(models.Friendship).DmChannelID, "\"}")

	fc.Response(w, http.StatusCreated, res)
}

func (fc *FriendshipController) HandleNewFriendshipRequest(w http.ResponseWriter, r *http.Request) {
	if fc.IsAllowedMethod(r, []string{"POST"}) == false {
		fc.Response(w, http.StatusMethodNotAllowed, nil)
		return
	}

	if fc.IsAuthorized(r) == false {
		fc.Response(w, http.StatusUnauthorized, nil)
		return
	}

	payload := fc.GetPayload(r)

	initiatorID := payload["initiator_id"].(string)
	friendID := payload["friend_id"].(string)

	if initiatorID == "" || friendID == "" {
		fc.Response(w, http.StatusBadRequest, nil)
		return
	}

	userID := strings.Split(strings.Split(r.Header.Get("Authorization"), "Bearer ")[1], ".")[0]

	if initiatorID != userID {
		fc.Response(w, http.StatusUnauthorized, nil)
		return
	}

	friendshipService := services.NewFriendshipService()
	newFriendshipRequest, err := friendshipService.CreateFriendshipRequest(fc.db, payload)

	if err != nil {
		fc.Response(w, http.StatusInternalServerError, nil)
		return
	}

	res := fmt.Sprint("{\"id\": \"", newFriendshipRequest.(models.FriendshipRequest).ID, "\"}")

	fc.Response(w, http.StatusCreated, res)
}

func (fc *FriendshipController) HandleGetFriendships(w http.ResponseWriter, r *http.Request) {
	if fc.IsAllowedMethod(r, []string{"GET"}) == false {
		fc.Response(w, http.StatusMethodNotAllowed, nil)
		return
	}

	if fc.IsAuthorized(r) == false {
		fc.Response(w, http.StatusUnauthorized, nil)
		return
	}

	userID := strings.Split(strings.Split(r.Header.Get("Authorization"), "Bearer ")[1], ".")[0]

	if userID == "" {
		fc.Response(w, http.StatusBadRequest, nil)
		return
	}

	friendshipService := services.NewFriendshipService()
	friendshipRecords, err := friendshipService.GetFriendships(fc.db, userID)

	if err != nil {
		fc.Response(w, http.StatusInternalServerError, nil)
		return
	}

	var res string
	for _, friendship := range friendshipRecords.([]models.Friendship) {
		res += fmt.Sprint("{\"id\": \"", friendship.ID, "\", \"channel_id\": \"", friendship.DmChannelID, "\"},")
	}

	fc.Response(w, http.StatusOK, fmt.Sprint("[", res[:len(res) - 1], "]"))
}

func (fc *FriendshipController) HandleGetFriendshipRequests(w http.ResponseWriter, r *http.Request) {
	if fc.IsAllowedMethod(r, []string{"GET"}) == false {
		fc.Response(w, http.StatusMethodNotAllowed, nil)
		return
	}

	if fc.IsAuthorized(r) == false {
		fc.Response(w, http.StatusUnauthorized, nil)
		return
	}

	userID := strings.Split(strings.Split(r.Header.Get("Authorization"), "Bearer ")[1], ".")[0]

	friendshipService := services.NewFriendshipService()
	friendshipRequestRecords, err := friendshipService.GetFriendshipRequests(fc.db, userID)

	if err != nil {
		fc.Response(w, http.StatusInternalServerError, nil)
		return
	}

	var res string
	for _, friendshipRequest := range friendshipRequestRecords.([]models.FriendshipRequest) {
		res += fmt.Sprint("{\"id\": \"", friendshipRequest.ID, "\", \"initiator_id\": \"", friendshipRequest.InitiatorID, "\", \"friend_id\": \"", friendshipRequest.FriendID, "\"},")
	}

	fc.Response(w, http.StatusOK, fmt.Sprint("[", res[:len(res) - 1], "]"))
}

func (fc *FriendshipController) HandleDeleteFriendship(w http.ResponseWriter, r *http.Request) {
	if fc.IsAllowedMethod(r, []string{"POST"}) == false {
		fc.Response(w, http.StatusMethodNotAllowed, nil)
		return
	}

	if fc.IsAuthorized(r) == false {
		fc.Response(w, http.StatusUnauthorized, nil)
		return
	}

	payload := fc.GetPayload(r)

	friendshipID := payload["id"].(string)

	if friendshipID == "" {
		fc.Response(w, http.StatusBadRequest, nil)
		return
	}

	friendshipService := services.NewFriendshipService()
	_, err := friendshipService.DeleteFriendship(fc.db, friendshipID)

	if err != nil {
		fc.Response(w, http.StatusInternalServerError, nil)
		return
	}

	fc.Response(w, http.StatusOK, nil)
}

func (fc *FriendshipController) HandleDeleteFriendshipRequest(w http.ResponseWriter, r *http.Request) {
	if fc.IsAllowedMethod(r, []string{"POST"}) == false {
		fc.Response(w, http.StatusMethodNotAllowed, nil)
		return
	}

	if fc.IsAuthorized(r) == false {
		fc.Response(w, http.StatusUnauthorized, nil)
		return
	}

	payload := fc.GetPayload(r)
	friendshipRequestID := payload["id"].(string)

	if friendshipRequestID == "" {
		fc.Response(w, http.StatusBadRequest, payload)
		return
	}

	friendshipService := services.NewFriendshipService()
	_, err := friendshipService.DeleteFriendshipRequest(fc.db, friendshipRequestID)
	
	if err != nil {
		fc.Response(w, http.StatusInternalServerError, payload)
		return
	}

	fc.Response(w, http.StatusOK, nil)
}

func NewFriendshipController(db *gorm.DB) *FriendshipController {
	friendshipController := &FriendshipController{
		db: db,
	}

	return &FriendshipController{
		Controller: *core.NewController([]core.ControllerMethod{
			friendshipController.HandleNewFriendship,
			friendshipController.HandleGetFriendships,
			friendshipController.HandleGetFriendshipRequests,
			friendshipController.HandleNewFriendshipRequest,
			friendshipController.HandleDeleteFriendship,
			friendshipController.HandleDeleteFriendshipRequest,
		}),
		db: db,
	}
}
