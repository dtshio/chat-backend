package services

import (
	"strconv"

	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"github.com/datsfilipe/pkg/repositories"
	"go.uber.org/zap"
)

type FriendshipService struct {
	core.Service
	log *zap.Logger
	repo *repositories.FriendshipRepository
	channelService *ChannelService
}

func (fs *FriendshipService) CreateFriendshipRequest(data interface{}) (interface{}, error) {
	payload := data.(core.Map)

	initiatorIDStr := payload["initiator_id"].(string)
	friendIDStr := payload["friend_id"].(string)

	initiatorID, err := strconv.Atoi(initiatorIDStr)
	if err != nil {
		fs.log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, fs.GenError(fs.InvalidData, initiatorIDStr)
	}

	friendID, err := strconv.Atoi(friendIDStr)
	if err != nil {
		fs.log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, fs.GenError(fs.InvalidData, friendIDStr)
	}

	friendshipRequest := &models.FriendshipRequest{}
	friendshipRequest.InitiatorID = models.BigInt(initiatorID)
	friendshipRequest.FriendID = models.BigInt(friendID)

	dbRecord, err := fs.repo.CreateFriendshipRequest(friendshipRequest)
	if err != nil {
		fs.log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, fs.GenError(fs.CreateError, dbRecord)
	}

	return dbRecord, nil
}

func (fs *FriendshipService) CreateFriendship(data interface{}) (interface{}, error) {
	payload := data.(core.Map)

	initiatorIDStr := payload["initiator_id"].(string)
	friendIDStr := payload["friend_id"].(string)
	requestID := core.HashID(initiatorIDStr, friendIDStr)

	initiatorID, err := strconv.Atoi(initiatorIDStr)
	if err != nil {
		fs.log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, fs.GenError(fs.InvalidData, nil)
	}

	friendID, err := strconv.Atoi(friendIDStr)
	if err != nil {
		fs.log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, fs.GenError(fs.InvalidData, nil)
	}

	_, err = fs.repo.AcceptFriendshipRequest(requestID)
	if err != nil {
		fs.log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	channel := &models.Channel{}
	channel.Type = "DIRECT_MESSAGE"

	channelRecord, err := fs.channelService.CreateChannel(channel)
	if err != nil {
		fs.log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	friendship := &models.Friendship{}
	friendship.ChannelID = channelRecord.(models.Channel).ID
	friendship.InitiatorID = models.BigInt(initiatorID)
	friendship.FriendID = models.BigInt(friendID)

	dbRecord, err := fs.repo.CreateFriendship(friendship)
	if err != nil {
		fs.log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecord, nil
}

func (fs *FriendshipService) GetFriendshipRequests(data interface{}) (interface{}, error) {
	userID := data.(string)

	if userID == "" {
		return nil, fs.GenError(fs.InvalidData, nil)
	}

	dbRecords, err := fs.repo.GetFriendshipRequests(userID)
	if err != nil {
		fs.log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecords, nil
}

func (fs *FriendshipService) GetFriendships(data interface{}) (interface{}, error) {
	userID, ok := data.(string)

	if !ok || userID == "" {
		return nil, fs.GenError(fs.InvalidData, nil)
	}

	dbRecords, err := fs.repo.GetFriendships(userID)
	if err != nil {
		fs.log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecords, nil
}

func (fs *FriendshipService) DeleteFriendshipRequest(data interface{}) (interface{}, error) {
	id := data.(string)

	_, err := fs.repo.DeleteFriendshipRequest(id)
	if err != nil {
		fs.log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return nil, nil
}

func (fs *FriendshipService) DeleteFriendship(data interface{}) (interface{}, error) {
	userID := data.(string)

	dbRecord, err := fs.repo.GetFriendship(userID)
	if err != nil {
		return nil, err
	}

	channelID := strconv.Itoa(int(dbRecord.(models.Friendship).ChannelID))

	_, err = fs.repo.DeleteFriendship(dbRecord)
	if err != nil {
		fs.log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	_, err = fs.channelService.DeleteChannel(channelID)
	if err != nil {
		fs.log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return nil, nil
}

func NewFriendshipService(
	log *zap.Logger,
	repo *repositories.FriendshipRepository,
	channelService *ChannelService,
) *FriendshipService {
	return &FriendshipService{
		Service: *core.NewService(
            []core.ServiceMethod{
                (&FriendshipService{}).CreateFriendship,
                (&FriendshipService{}).GetFriendships,
                (&FriendshipService{}).DeleteFriendship,
                (&FriendshipService{}).CreateFriendshipRequest,
                (&FriendshipService{}).GetFriendshipRequests,
                (&FriendshipService{}).DeleteFriendshipRequest,
            },
		),
		log:            log,
		repo:           repo,
		channelService: channelService,
	}
}
