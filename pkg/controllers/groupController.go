package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"github.com/datsfilipe/pkg/services"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type GroupController struct {
	core.Controller
	db *gorm.DB
	log *zap.Logger
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

	service := services.NewGroupService()

	groupRecord, err := service.CreateGroup(mc.db, mc.log, group)
	if err != nil {
		mc.Response(w, http.StatusInternalServerError, err)
		return
	}

	groupMember := &models.GroupMember{
		GroupID: groupRecord.(models.Group).ID,
		ProfileID: models.BigInt(profileID),
		ChannelID: groupRecord.(models.Group).ChannelID,
	}

	memberRecord, err := service.AddGroupMember(mc.db, mc.log, groupMember)
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

	service := services.NewGroupService()

	dbRecords, err := service.GetGroups(mc.db, mc.log, groups)
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
	if mc.IsAllowedMethod(r, []string{"POST"}) == false {
		mc.Response(w, http.StatusMethodNotAllowed, nil)
		return
	}

	payload := mc.GetPayload(r)

	if payload == nil {
		mc.Response(w, http.StatusBadRequest, nil)
		return
	}

	id, ok := payload["id"].(string)

	if !ok || id == "" {
		mc.Response(w, http.StatusBadRequest, nil)
		return
	}

	service := services.NewGroupService()

	_, err := service.DeleteGroup(mc.db, mc.log, id)
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

	service := services.NewGroupService()

	dbRecord, err := service.AddGroupMember(mc.db, mc.log, groupMember)
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

	service := services.NewGroupService()

	dbRecords, err := service.GetGroupMembers(mc.db, mc.log, groupID)
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
	if mc.IsAllowedMethod(r, []string{"POST"}) == false {
		mc.Response(w, http.StatusMethodNotAllowed, nil)
		return
	}

	payload := mc.GetPayload(r)

	if payload == nil {
		mc.Response(w, http.StatusBadRequest, nil)
		return
	}

	id, ok := payload["id"].(string)

	if !ok || id == "" {
		mc.Response(w, http.StatusBadRequest, nil)
		return
	}

	service := services.NewGroupService()

	_, err := service.DeleteGroupMember(mc.db, mc.log, id)
	if err != nil {
		mc.Response(w, http.StatusInternalServerError, err)
		return
	}

	mc.Response(w, http.StatusOK, "Deleted")
}

func (mc *GroupController) HandleGetMember(w http.ResponseWriter, r *http.Request) {
	if mc.IsAllowedMethod(r, []string{"POST"}) == false {
		mc.Response(w, http.StatusMethodNotAllowed, nil)
		return
	}

	payload := mc.GetPayload(r)

	if payload == nil {
		mc.Response(w, http.StatusBadRequest, nil)
		return
	}

	id, ok := payload["id"].(string)

	if !ok || id == "" {
		mc.Response(w, http.StatusBadRequest, nil)
		return
	}

	service := services.NewGroupService()

	dbRecord, err := service.GetGroupMember(mc.db, mc.log, id)
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
	if mc.IsAllowedMethod(r, []string{"POST"}) == false {
		mc.Response(w, http.StatusMethodNotAllowed, nil)
		return
	}

	payload := mc.GetPayload(r)
	
	if payload == nil {
		mc.Response(w, http.StatusBadRequest, nil)
		return
	}

	profileID := payload["profile_id"].(string)

	service := services.NewGroupService()

	dbRecords, err := service.GetGroupsByProfile(mc.db, mc.log, profileID)
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

func NewGroupController(db *gorm.DB, log *zap.Logger) *GroupController {
	controller := &GroupController{
		db: db,
		log: log,
	}

	return &GroupController{
		Controller: *core.NewController([]core.ControllerMethod{
			controller.HandleNewGroup,
			controller.HandleGetGroups,
			controller.HandleDeleteGroup,
			controller.HandleAddMember,
			controller.HandleGetMembers,
			controller.HandleRemoveMember,
			controller.HandleGetMember,
			controller.HandleGetGroupsByProfile,
		}),
		db: db,
		log: log,
	}
}
