package services

import (
	"strconv"

	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"github.com/datsfilipe/pkg/repositories"
	"go.uber.org/zap"
)

type GroupService struct {
	core.Service
	log *zap.Logger
	repo *repositories.GroupRepository
	service *ChannelService
}

func (gs *GroupService) CreateGroup(data interface{}) (interface{}, error) {
	group := data.(*models.Group)

	channel := &models.Channel{}
	channel.Type = "GROUP"

	channelRecord, err := gs.service.CreateChannel(channel)
	if err != nil {
		gs.log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	group.ChannelID = channelRecord.(models.Channel).ID

	dbRecord, err := gs.repo.CreateGroup(group)
	if err != nil {
		gs.log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecord, nil
}

func (gs *GroupService) GetGroup(data interface{}) (interface{}, error) {
	id := data.(string)

	dbRecord, err := gs.repo.GetGroup(id)
	if err != nil {
		gs.log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecord, nil
}

func (gs *GroupService) GetGroups(data interface{}) (interface{}, error) {
	groups := data.(*[]models.Group)

	if groups == nil {
		return nil, gs.GenError(gs.InvalidData, groups)
	}

	dbRecords, err := gs.repo.GetGroups(groups)
	if err != nil {
		gs.log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecords, nil
}

func (gs *GroupService) DeleteGroup(data interface{}) (interface{}, error) {
	id, ok := data.(string)

	if !ok || id == "" {
		return nil, gs.GenError(gs.InvalidData, nil)
	}

	dbRecord, err := gs.repo.DeleteGroup(id)
	if err != nil {
		gs.log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecord, nil
}

func (gs *GroupService) AddGroupMember(data interface{}) (interface{}, error) {
	groupMember := data.(*models.GroupMember)
	profileID := strconv.FormatInt(int64(groupMember.ProfileID), 10)

	profileGroups, err := gs.repo.GetGroupsByProfile(profileID)
	if err != nil {
		gs.log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	for _, profileGroup := range profileGroups.([]models.Group) {
		if profileGroup.ID == groupMember.GroupID {
			return nil, gs.GenError(gs.DuplicateError, nil)
		}
	}

	groupID := strconv.FormatInt(int64(groupMember.GroupID), 10)

	groupRecord, err := gs.repo.GetGroup(groupID)
	if err != nil {
		gs.log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	groupMember.ChannelID = groupRecord.(models.Group).ChannelID

	dbRecord, err := gs.repo.AddGroupMember(groupMember)
	if err != nil {
		gs.log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecord, nil
}

func (gs *GroupService) GetGroupMembers(data interface{}) (interface{}, error) {
	groupID := data.(string)

	dbRecords, err := gs.repo.GetGroupMembers(groupID)
	if err != nil {
		gs.log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecords, nil
}

func (gs *GroupService) DeleteGroupMember(data interface{}) (interface{}, error) {
	groupMemberID := data.(string)

	dbRecord, err := gs.repo.DeleteGroupMember(groupMemberID)
	if err != nil {
		gs.log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecord, nil
}

func (gs *GroupService) GetGroupMember(data interface{}) (interface{}, error) {
	groupMemberID := data.(string)

	dbRecord, err := gs.repo.GetGroupMember(groupMemberID)
	if err != nil {
		gs.log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecord, nil
}

func (gs *GroupService) GetGroupsByProfile(data interface{}) (interface{}, error) {
	profileID, ok := data.(string)

	if !ok || profileID == "" {
		return nil, gs.GenError(gs.InvalidData, nil)
	}

	dbRecords, err := gs.repo.GetGroupsByProfile(profileID)
	if err != nil {
		gs.log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecords, nil
}

func NewGroupService(log *zap.Logger, repo *repositories.GroupRepository, channelService *ChannelService) *GroupService {
	return &GroupService{
		Service: *core.NewService(),
		log: log,
		repo: repo,
		service: channelService,
	}
}
