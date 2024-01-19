package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"github.com/datsfilipe/pkg/services"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type FriendshipController struct {
	core.Controller
	db *gorm.DB
	log *zap.Logger
}

func (fc *FriendshipController) HandleNewFriendship(w http.ResponseWriter, r *http.Request) {
	if fc.IsAllowedMethod(r, []string{"POST"}) == false {
		fc.Response(w, http.StatusMethodNotAllowed, nil)
		return
	}

	payload := fc.GetPayload(r)

	initiatorID := payload["initiator_id"].(string)
	friendID := payload["friend_id"].(string)

	if initiatorID == "" || friendID == "" {
		fc.Response(w, http.StatusBadRequest, nil)
		return
	}

	service := services.NewFriendshipService()

	dbRecord, err := service.CreateFriendship(fc.db, fc.log, payload)
	if err != nil {
		fc.Response(w, http.StatusInternalServerError, err)
		return
	}

	res := fmt.Sprint("{\"id\": \"", dbRecord.(models.Friendship).ID, ", \"channel_id\": \"", dbRecord.(models.Friendship).DmChannelID, "\"}")

	fc.Response(w, http.StatusCreated, res)
}

func (fc *FriendshipController) HandleNewFriendshipRequest(w http.ResponseWriter, r *http.Request) {
	if fc.IsAllowedMethod(r, []string{"POST"}) == false {
		fc.Response(w, http.StatusMethodNotAllowed, nil)
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

	service := services.NewFriendshipService()

	dbRecord, err := service.CreateFriendshipRequest(fc.db, fc.log, payload)
	if err != nil {
		fc.Response(w, http.StatusInternalServerError, err)
		return
	}

	res := fmt.Sprint("{\"id\": \"", dbRecord.(models.FriendshipRequest).ID, "\"}")

	fc.Response(w, http.StatusCreated, res)
}

func (fc *FriendshipController) HandleGetFriendships(w http.ResponseWriter, r *http.Request) {
	if fc.IsAllowedMethod(r, []string{"GET"}) == false {
		fc.Response(w, http.StatusMethodNotAllowed, nil)
		return
	}

	userID := strings.Split(strings.Split(r.Header.Get("Authorization"), "Bearer ")[1], ".")[0]

	if userID == "" {
		fc.Response(w, http.StatusBadRequest, nil)
		return
	}

	service := services.NewFriendshipService()

	dbRecords, err := service.GetFriendships(fc.db, fc.log, userID)
	if err != nil {
		fc.Response(w, http.StatusInternalServerError, err)
		return
	}

	if dbRecords == nil {
		fc.Response(w, http.StatusOK, "[]")
		return
	}

	var res string
	for _, friendship := range dbRecords.([]models.Friendship) {
		res += fmt.Sprint("{\"id\": \"", friendship.ID, "\", \"initiator_id\": \"", friendship.InitiatorID, "\", \"friend_id\": \"", friendship.FriendID, "\"},")
	}
	fc.log.Info(res)

	fc.Response(w, http.StatusOK, fmt.Sprint("[", res[:len(res) - 1], "]"))
}

func (fc *FriendshipController) HandleGetFriendshipRequests(w http.ResponseWriter, r *http.Request) {
	if fc.IsAllowedMethod(r, []string{"GET"}) == false {
		fc.Response(w, http.StatusMethodNotAllowed, nil)
		return
	}

	userID := strings.Split(strings.Split(r.Header.Get("Authorization"), "Bearer ")[1], ".")[0]

	service := services.NewFriendshipService()

	dbRecords, err := service.GetFriendshipRequests(fc.db, fc.log, userID)
	if err != nil {
		fc.Response(w, http.StatusInternalServerError, err)
		return
	}

	var res string
	for _, friendshipRequest := range dbRecords.([]models.FriendshipRequest) {
		res += fmt.Sprint("{\"id\": \"", friendshipRequest.ID, "\", \"initiator_id\": \"", friendshipRequest.InitiatorID, "\", \"friend_id\": \"", friendshipRequest.FriendID, "\"},")
	}

	fc.Response(w, http.StatusOK, fmt.Sprint("[", res[:len(res) - 1], "]"))
}

func (fc *FriendshipController) HandleDeleteFriendship(w http.ResponseWriter, r *http.Request) {
	if fc.IsAllowedMethod(r, []string{"POST"}) == false {
		fc.Response(w, http.StatusMethodNotAllowed, nil)
		return
	}

	payload := fc.GetPayload(r)

	friendshipID := payload["id"].(string)

	if friendshipID == "" {
		fc.Response(w, http.StatusBadRequest, nil)
		return
	}

	service := services.NewFriendshipService()

	_, err := service.DeleteFriendship(fc.db, fc.log, friendshipID)
	if err != nil {
		fc.Response(w, http.StatusInternalServerError, err)
		return
	}

	fc.Response(w, http.StatusOK, "Deleted")
}

func (fc *FriendshipController) HandleDeleteFriendshipRequest(w http.ResponseWriter, r *http.Request) {
	if fc.IsAllowedMethod(r, []string{"POST"}) == false {
		fc.Response(w, http.StatusMethodNotAllowed, nil)
		return
	}

	payload := fc.GetPayload(r)
	requestID := payload["id"].(string)

	if requestID == "" {
		fc.Response(w, http.StatusBadRequest, payload)
		return
	}

	service := services.NewFriendshipService()

	_, err := service.DeleteFriendshipRequest(fc.db, fc.log, requestID)
	if err != nil {
		fc.Response(w, http.StatusInternalServerError, err)
		return
	}

	fc.Response(w, http.StatusOK, "Deleted")
}

func NewFriendshipController(db *gorm.DB, log *zap.Logger) *FriendshipController {
	controller := &FriendshipController{
		db: db,
		log: log,
	}

	return &FriendshipController{
		Controller: *core.NewController([]core.ControllerMethod{
			controller.HandleNewFriendship,
			controller.HandleGetFriendships,
			controller.HandleGetFriendshipRequests,
			controller.HandleNewFriendshipRequest,
			controller.HandleDeleteFriendship,
			controller.HandleDeleteFriendshipRequest,
		}),
		db: db,
		log: log,
	}
}
