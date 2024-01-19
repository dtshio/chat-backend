package services

import (
	"strconv"

	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"github.com/datsfilipe/pkg/repositories"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type FriendshipService struct {
	core.Service
}

func (fs *FriendshipService) CreateFriendshipRequest(db *gorm.DB, log *zap.Logger, data interface{}) (interface{}, error) {
	payload := data.(core.Map)

	initiatorIDStr := payload["initiator_id"].(string)
	friendIDStr := payload["friend_id"].(string)

	initiatorID, err := strconv.Atoi(initiatorIDStr)
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, fs.GenError(fs.InvalidData, initiatorIDStr)
	}

	friendID, err := strconv.Atoi(friendIDStr)
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, fs.GenError(fs.InvalidData, friendIDStr)
	}

	repo := repositories.NewFriendshipRepository()

	friendshipRequest := &models.FriendshipRequest{}
	friendshipRequest.InitiatorID = models.BigInt(initiatorID)
	friendshipRequest.FriendID = models.BigInt(friendID)

	dbRecord, err := repo.CreateFriendshipRequest(db, friendshipRequest)
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, fs.GenError(fs.CreateError, dbRecord)
	}

	return dbRecord, nil
}

func (fs *FriendshipService) CreateFriendship(db *gorm.DB, log *zap.Logger, data interface{}) (interface{}, error) {
	payload := data.(core.Map)

	initiatorIDStr := payload["initiator_id"].(string)
	friendIDStr := payload["friend_id"].(string)
	requestID := core.HashID(initiatorIDStr, friendIDStr)

	initiatorID, err := strconv.Atoi(initiatorIDStr)
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, fs.GenError(fs.InvalidData, nil)
	}

	friendID, err := strconv.Atoi(friendIDStr)
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, fs.GenError(fs.InvalidData, nil)
	}

	repo := repositories.NewFriendshipRepository()

	_, err = repo.AcceptFriendshipRequest(db, requestID)
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	channelService := NewChannelService()

	channel, err := channelService.CreateChannel(db, log, &models.Channel{})
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	friendship := &models.Friendship{}
	friendship.DmChannelID = channel.(models.Channel).ID
	friendship.InitiatorID = models.BigInt(initiatorID)
	friendship.FriendID = models.BigInt(friendID)

	dbRecord, err := repo.CreateFriendship(db, friendship)
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecord, nil
}

func (fs *FriendshipService) GetFriendshipRequests(db *gorm.DB, log *zap.Logger, data interface{}) (interface{}, error) {
	userID := data.(string)

	if userID == "" {
		return nil, fs.GenError(fs.InvalidData, nil)
	}

	repo := repositories.NewFriendshipRepository()

	dbRecords, err := repo.GetFriendshipRequests(db, userID)
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecords, nil
}

func (fs *FriendshipService) GetFriendships(db *gorm.DB, log *zap.Logger, data interface{}) (interface{}, error) {
	userID, ok := data.(string)

	if !ok || userID == "" {
		return nil, fs.GenError(fs.InvalidData, nil)
	}

	repo := repositories.NewFriendshipRepository()

	dbRecords, err := repo.GetFriendships(db, userID)
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecords, nil
}

func (fs *FriendshipService) DeleteFriendshipRequest(db *gorm.DB, log *zap.Logger, data interface{}) (interface{}, error) {
	id := data.(string)

	repo := repositories.NewFriendshipRepository()

	_, err := repo.DeleteFriendshipRequest(db, id)
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return nil, nil
}

func (fs *FriendshipService) DeleteFriendship(db *gorm.DB, log *zap.Logger, data interface{}) (interface{}, error) {
	userID := data.(string)

	repo := repositories.NewFriendshipRepository()

	dbRecord, err := repo.GetFriendship(db, userID)
	if err != nil {
		return nil, err
	}

	channelID := strconv.Itoa(int(dbRecord.(models.Friendship).DmChannelID))

	_, err = repo.DeleteFriendship(db, dbRecord)
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	channelService := NewChannelService()

	_, err = channelService.DeleteChannel(db, log, channelID)
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return nil, nil
}

func NewFriendshipService() *FriendshipService {
	service := &FriendshipService{}

	return &FriendshipService{
		Service: *core.NewService(
			[]core.ServiceMethod{
				service.CreateFriendship,
				service.GetFriendships,
				service.DeleteFriendship,
				service.CreateFriendshipRequest,
				service.GetFriendshipRequests,
				service.DeleteFriendshipRequest,
			},
		),
	}
}
