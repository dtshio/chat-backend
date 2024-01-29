package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"github.com/datsfilipe/pkg/services"
)

type FriendshipController struct {
	core.Controller
	service *services.FriendshipService
	userService *services.UserService
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

	dbRecord, err := fc.service.CreateFriendship(payload)
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

	dbRecord, err := fc.service.CreateFriendshipRequest(payload)
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

	dbRecords, err := fc.service.GetFriendships(userID)
	if err != nil {
		if strings.Contains(err.Error(), "Entry not found") {
			fc.Response(w, http.StatusOK, "[]")
			return
		}

		fc.Response(w, http.StatusInternalServerError, err)
		return
	}

	if dbRecords == nil {
		fc.Response(w, http.StatusOK, "[]")
		return
	}

	var res string
	for _, friendship := range dbRecords.([]models.Friendship) {
		res += fmt.Sprint(
			"{\"id\": \"",
			friendship.ID,
			"\", \"initiator_id\": \"",
			friendship.InitiatorID,
			"\", \"friend_id\": \"",
			friendship.FriendID,
			"\", \"channel_id\": \"",
			friendship.DmChannelID,
			"\"},",
		)
	}

	fc.Response(w, http.StatusOK, fmt.Sprint("[", res[:len(res) - 1], "]"))
}

func (fc *FriendshipController) HandleGetFriendshipRequests(w http.ResponseWriter, r *http.Request) {
	if fc.IsAllowedMethod(r, []string{"GET"}) == false {
		fc.Response(w, http.StatusMethodNotAllowed, nil)
		return
	}

	userID := strings.Split(strings.Split(r.Header.Get("Authorization"), "Bearer ")[1], ".")[0]

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		fc.Response(w, http.StatusInternalServerError, err)
		return
	}

	friendshipRequestRecords, err := fc.service.GetFriendshipRequests(userID)
	if err != nil {
		if strings.Contains(err.Error(), "Entry not found") {
			fc.Response(w, http.StatusOK, "[]")
			return
		}

		fc.Response(w, http.StatusInternalServerError, err)
		return
	}

	var friendProfiles []models.Profile
	var userIDs []models.BigInt
	for _, friendshipRequest := range friendshipRequestRecords.([]models.FriendshipRequest) {
		if friendshipRequest.InitiatorID == models.BigInt(userIDInt) {
			userIDs = append(userIDs, friendshipRequest.FriendID)
		} else {
			userIDs = append(userIDs, friendshipRequest.InitiatorID)
		}
	}

	var userIDsStr []string
	for _, userID := range userIDs {
		userIDsStr = append(userIDsStr, fmt.Sprint(userID))
	}

	if profilesRecords, err := fc.userService.GetDefaultProfiles(userIDsStr); err != nil {
		fc.Response(w, http.StatusInternalServerError, err)
		return
	} else {
		friendProfiles = profilesRecords.([]models.Profile)
	}

	var res []core.Map

	for _, friendshipRequest := range friendshipRequestRecords.([]models.FriendshipRequest) {
		var profile models.Profile

		for _, friendProfile := range friendProfiles {
			if friendshipRequest.InitiatorID == friendProfile.UserID {
				profile = friendProfile
				break
			}
		}

		if profile.ID == 0 {
			break
		}

		res = append(res, core.Map{
			"id": profile.ID,
			"request_id": friendshipRequest.ID,
			"username": profile.Username,
			"initiator_id": friendshipRequest.InitiatorID,
			"friend_id": friendshipRequest.FriendID,
		})
	}

	if len(res) == 0 {
		fc.Response(w, http.StatusOK, "[]")
		return
	}

	fc.Response(w, http.StatusOK, res)
}

func (fc *FriendshipController) HandleDeleteFriendship(w http.ResponseWriter, r *http.Request) {
	if fc.IsAllowedMethod(r, []string{"DELETE"}) == false {
		fc.Response(w, http.StatusMethodNotAllowed, nil)
		return
	}

	friendshipID := fc.GetPayload(r)["id"].(string)

	if friendshipID == "" {
		fc.Response(w, http.StatusBadRequest, nil)
		return
	}

	_, err := fc.service.DeleteFriendship(friendshipID)
	if err != nil {
		fc.Response(w, http.StatusInternalServerError, err)
		return
	}

	fc.Response(w, http.StatusOK, "Deleted")
}

func (fc *FriendshipController) HandleDeleteFriendshipRequest(w http.ResponseWriter, r *http.Request) {
	if fc.IsAllowedMethod(r, []string{"DELETE"}) == false {
		fc.Response(w, http.StatusMethodNotAllowed, nil)
		return
	}

	requestID := fc.GetPayload(r)["id"].(string)

	if requestID == "" {
		fc.Response(w, http.StatusBadRequest, requestID)
		return
	}

	_, err := fc.service.DeleteFriendshipRequest(requestID)
	if err != nil {
		fc.Response(w, http.StatusInternalServerError, err)
		return
	}

	fc.Response(w, http.StatusOK, "Deleted")
}

func NewFriendshipController(
	service *services.FriendshipService,
	userService *services.UserService,
) *FriendshipController {
	return &FriendshipController{
		service: service,
		userService: userService,
	}
}
