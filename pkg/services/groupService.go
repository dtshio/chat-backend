package services

import (
	"strconv"

	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"github.com/datsfilipe/pkg/repositories"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type GroupService struct {
	core.Service
}

func (gs *GroupService) CreateGroup(db *gorm.DB, log *zap.Logger, data interface{}) (interface{}, error) {
	group := data.(*models.Group)

	repo := repositories.NewGroupRepository()

	channelService := NewChannelService()

	channel := &models.Channel{}
	channel.Type = "GROUP"

	channelRecord, err := channelService.CreateChannel(db, log, channel)
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	group.ChannelID = channelRecord.(models.Channel).ID

	dbRecord, err := repo.CreateGroup(db, group)
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecord, nil
}

func (gs *GroupService) GetGroups(db *gorm.DB, log *zap.Logger, data interface{}) (interface{}, error) {
	groups := data.(*[]models.Group)

	if groups == nil {
		return nil, gs.GenError(gs.InvalidData, groups)
	}

	repo := repositories.NewGroupRepository()

	dbRecords, err := repo.GetGroups(db, groups)
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecords, nil
}

func (gs *GroupService) DeleteGroup(db *gorm.DB, log *zap.Logger, data interface{}) (interface{}, error) {
	id, ok := data.(string)

	if !ok || id == "" {
		return nil, gs.GenError(gs.InvalidData, nil)
	}

	repo := repositories.NewGroupRepository()

	dbRecord, err := repo.DeleteGroup(db, id)
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecord, nil
}

func (gs *GroupService) AddGroupMember(db *gorm.DB, log *zap.Logger, data interface{}) (interface{}, error) {
	groupMember := data.(*models.GroupMember)
	profileID := strconv.FormatInt(int64(groupMember.ProfileID), 10)

	repo := repositories.NewGroupRepository()

	profileGroups, err := repo.GetGroupsByProfile(db, profileID)
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	for _, profileGroup := range profileGroups.([]models.Group) {
		if profileGroup.ID == groupMember.GroupID {
			return nil, gs.GenError(gs.DuplicateError, nil)
		}
	}

	groupID := strconv.FormatInt(int64(groupMember.GroupID), 10)

	groupRecord, err := repo.GetGroup(db, groupID)
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	groupMember.ChannelID = groupRecord.(models.Group).ChannelID

	dbRecord, err := repo.AddGroupMember(db, groupMember)
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecord, nil
}

func (gs *GroupService) GetGroupMembers(db *gorm.DB, log *zap.Logger, data interface{}) (interface{}, error) {
	groupID := data.(string)

	repo := repositories.NewGroupRepository()

	dbRecords, err := repo.GetGroupMembers(db, groupID)
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecords, nil
}

func (gs *GroupService) DeleteGroupMember(db *gorm.DB, log *zap.Logger, data interface{}) (interface{}, error) {
	groupMemberID := data.(string)

	repo := repositories.NewGroupRepository()

	dbRecord, err := repo.DeleteGroupMember(db, groupMemberID)
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecord, nil
}

func (gs *GroupService) GetGroupMember(db *gorm.DB, log *zap.Logger, data interface{}) (interface{}, error) {
	groupMemberID := data.(string)

	repo := repositories.NewGroupRepository()

	dbRecord, err := repo.GetGroupMember(db, groupMemberID)
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecord, nil
}

func (gs *GroupService) GetGroupsByProfile(db *gorm.DB, log *zap.Logger, data interface{}) (interface{}, error) {
	profileID, ok := data.(string)

	if !ok || profileID == "" {
		return nil, gs.GenError(gs.InvalidData, nil)
	}

	repo := repositories.NewGroupRepository()

	dbRecords, err := repo.GetGroupsByProfile(db, profileID)
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecords, nil
}

func NewGroupService() *GroupService {
	service := &GroupService{}

	return &GroupService{
		Service: *core.NewService(
			[]core.ServiceMethod{
				service.CreateGroup,
				service.GetGroups,
				service.DeleteGroup,
				service.AddGroupMember,
				service.GetGroupMembers,
				service.DeleteGroupMember,
				service.GetGroupMember,
				service.GetGroupsByProfile,
			},
		),
	}
}
