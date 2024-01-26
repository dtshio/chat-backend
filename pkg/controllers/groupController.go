package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"github.com/datsfilipe/pkg/services"
)

type GroupController struct {
	core.Controller
	service *services.GroupService
}

func (mc *GroupController) HandleNewGroup(w http.ResponseWriter, r *http.Request) {
	if mc.IsAllowedMethod(r, []string{"POST"}) == false {
		mc.Response(w, http.StatusMethodNotAllowed, nil)
		return
	}

	payload := mc.GetPayload(r)

	if payload == nil {
		mc.Response(w, http.StatusBadRequest, nil)
		return
	}

	group := &models.Group{}

	ownerID, _ := strconv.ParseInt(payload["owner_id"].(string), 10, 64)

	group.Name = payload["name"].(string)
	group.OwnerID = models.BigInt(ownerID)
	profileID, _ := strconv.ParseInt(payload["profile_id"].(string), 10, 64)

	if group.Name == "" || group.OwnerID == 0 || profileID == 0 {
		mc.Response(w, http.StatusBadRequest, nil)
		return
	}

	groupRecord, err := mc.service.CreateGroup(group)
	if err != nil {
		mc.Response(w, http.StatusInternalServerError, err)
		return
	}

	groupMember := &models.GroupMember{
		GroupID: groupRecord.(models.Group).ID,
		ProfileID: models.BigInt(profileID),
		ChannelID: groupRecord.(models.Group).ChannelID,
	}

	memberRecord, err := mc.service.AddGroupMember(groupMember)
	if err != nil {
		mc.Response(w, http.StatusInternalServerError, err)
		return
	}

	res, err := json.Marshal(memberRecord)
	if err != nil {
		mc.Response(w, http.StatusInternalServerError, nil)
		return
	}

	mc.Response(w, http.StatusCreated, res)
}

func (mc *GroupController) HandleGetGroups(w http.ResponseWriter, r *http.Request) {
	if mc.IsAllowedMethod(r, []string{"POST"}) == false {
		mc.Response(w, http.StatusMethodNotAllowed, nil)
		return
	}

	payload := mc.GetPayload(r)

	if payload == nil {
		mc.Response(w, http.StatusBadRequest, nil)
		return
	}

	groups := &[]models.Group{}

	dbRecords, err := mc.service.GetGroups(groups)
	if err != nil {
		mc.Response(w, http.StatusInternalServerError, err)
		return
	}

	res, err := json.Marshal(dbRecords)
	if err != nil {
		mc.Response(w, http.StatusInternalServerError, nil)
		return
	}

	mc.Response(w, http.StatusOK, res)
}

func (mc *GroupController) HandleDeleteGroup(w http.ResponseWriter, r *http.Request) {
	if mc.IsAllowedMethod(r, []string{"DELETE"}) == false {
		mc.Response(w, http.StatusMethodNotAllowed, nil)
		return
	}

	id := mc.GetPayload(r)["id"].(string)

	if id == "" {
		mc.Response(w, http.StatusBadRequest, nil)
		return
	}

	_, err := mc.service.DeleteGroup(id)
	if err != nil {
		mc.Response(w, http.StatusInternalServerError, err)
		return
	}

	mc.Response(w, http.StatusOK, "Deleted")
}

func (mc *GroupController) HandleAddMember(w http.ResponseWriter, r *http.Request) {
	if mc.IsAllowedMethod(r, []string{"POST"}) == false {
		mc.Response(w, http.StatusMethodNotAllowed, nil)
		return
	}

	payload := mc.GetPayload(r)

	groupMember := &models.GroupMember{}

	if payload == nil {
		mc.Response(w, http.StatusBadRequest, nil)
		return
	}

	groupID, _ := strconv.ParseInt(payload["group_id"].(string), 10, 64)
	profileID, _ := strconv.ParseInt(payload["profile_id"].(string), 10, 64)

	groupMember.GroupID = models.BigInt(groupID)
	groupMember.ProfileID = models.BigInt(profileID)

	if groupMember.GroupID == 0 || groupMember.ProfileID == 0 {
		mc.Response(w, http.StatusBadRequest, nil)
		return
	}

	dbRecord, err := mc.service.AddGroupMember(groupMember)
	if err != nil {
		mc.Response(w, http.StatusInternalServerError, err)
		return
	}

	res, err := json.Marshal(dbRecord)
	if err != nil {
		mc.Response(w, http.StatusInternalServerError, nil)
		return
	}

	mc.Response(w, http.StatusCreated, res)
}

func (mc *GroupController) HandleGetMembers(w http.ResponseWriter, r *http.Request) {
	if mc.IsAllowedMethod(r, []string{"POST"}) == false {
		mc.Response(w, http.StatusMethodNotAllowed, nil)
		return
	}

	payload := mc.GetPayload(r)

	if payload == nil {
		mc.Response(w, http.StatusBadRequest, nil)
		return
	}

	groupID, ok := payload["id"].(string)

	if !ok || groupID == "" {
		mc.Response(w, http.StatusBadRequest, nil)
		return
	}

	dbRecords, err := mc.service.GetGroupMembers(groupID)
	if err != nil {
		mc.Response(w, http.StatusInternalServerError, err)
		return
	}

	res, err := json.Marshal(dbRecords)
	if err != nil {
		mc.Response(w, http.StatusInternalServerError, nil)
		return
	}

	mc.Response(w, http.StatusOK, res)
}

func (mc *GroupController) HandleRemoveMember(w http.ResponseWriter, r *http.Request) {
	if mc.IsAllowedMethod(r, []string{"DELETE"}) == false {
		mc.Response(w, http.StatusMethodNotAllowed, nil)
		return
	}

	id := mc.GetPayload(r)["id"].(string)

	if id == "" {
		mc.Response(w, http.StatusBadRequest, nil)
		return
	}

	_, err := mc.service.DeleteGroupMember(id)
	if err != nil {
		mc.Response(w, http.StatusInternalServerError, err)
		return
	}

	mc.Response(w, http.StatusOK, "Deleted")
}

func (mc *GroupController) HandleGetMember(w http.ResponseWriter, r *http.Request) {
	if mc.IsAllowedMethod(r, []string{"GET"}) == false {
		mc.Response(w, http.StatusMethodNotAllowed, nil)
		return
	}

	id := mc.GetPayload(r)["id"].(string)

	if id == "" {
		mc.Response(w, http.StatusBadRequest, nil)
		return
	}

	dbRecord, err := mc.service.GetGroupMember(id)
	if err != nil {
		mc.Response(w, http.StatusInternalServerError, err)
		return
	}

	res, err := json.Marshal(dbRecord)
	if err != nil {
		mc.Response(w, http.StatusInternalServerError, nil)
		return
	}

	mc.Response(w, http.StatusOK, res)
}

func (mc *GroupController) HandleGetGroupsByProfile(w http.ResponseWriter, r *http.Request) {
	if mc.IsAllowedMethod(r, []string{"GET"}) == false {
		mc.Response(w, http.StatusMethodNotAllowed, nil)
		return
	}

	payload := mc.GetPayload(r)
	profileID := payload["id"].(string)

	dbRecords, err := mc.service.GetGroupsByProfile(profileID)
	if err != nil {
		mc.Response(w, http.StatusInternalServerError, err)
		return
	}

	res, err := json.Marshal(dbRecords)
	if err != nil {
		mc.Response(w, http.StatusInternalServerError, nil)
		return
	}

	mc.Response(w, http.StatusOK, res)
}

func NewGroupController(service *services.GroupService) *GroupController {
	return &GroupController{
		service: service,
	}
}
