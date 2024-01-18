package services

import (
	"strconv"

	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"github.com/datsfilipe/pkg/repositories"
	"gorm.io/gorm"
)

type FriendshipService struct {
	core.Service
}

func (fs *FriendshipService) CreateFriendshipRequest(db *gorm.DB, data interface{}) (interface{}, error) {
	payload := data.(core.Map)

	initiatorID := payload["initiator_id"].(string)
	friendID := payload["friend_id"].(string)

	initiatorIDInt, err := strconv.Atoi(initiatorID)
	if err != nil {
		return nil, fs.GenError(fs.InvalidData, nil)
	}

	friendIDInt, err := strconv.Atoi(friendID)
	if err != nil {
		return nil, fs.GenError(fs.InvalidData, nil)
	}

	friendshipRepo := repositories.NewFriendshipRepository()

	friendshipRequest := &models.FriendshipRequest{}

	friendshipRequest.InitiatorID = models.BigInt(initiatorIDInt)
	friendshipRequest.FriendID = models.BigInt(friendIDInt)

	friendshipRequest.BeforeCreateRecord()

	friendshipRequestRecord, err := friendshipRepo.CreateFriendshipRequest(db, friendshipRequest)
	if err != nil {
		return nil, fs.GenError(fs.CreatingError, friendshipRequest)
	}

	return friendshipRequestRecord, nil
}

func (fs *FriendshipService) CreateFriendship(db *gorm.DB, data interface{}) (interface{}, error) {
	payload := data.(core.Map)

	initiatorID := payload["initiator_id"].(string)
	friendID := payload["friend_id"].(string)
	requestID := core.HashID(initiatorID, friendID)

	initiatorIDInt, err := strconv.Atoi(initiatorID)
	if err != nil {
		return nil, fs.GenError(fs.InvalidData, nil)
	}

	friendIDInt, err2 := strconv.Atoi(friendID)
	if err2 != nil {
		return nil, fs.GenError(fs.InvalidData, nil)
	}

	friendshipRepo := repositories.NewFriendshipRepository()

	_, err3 := friendshipRepo.AcceptFriendshipRequest(db, requestID)
	if err3 != nil {
		return nil, fs.GenError(fs.UpdateError, nil)
	}

	channelService := NewChannelService()

	channel, err4 := channelService.CreateChannel(db, &models.Channel{})
	if err4 != nil {
		return nil, fs.GenError(fs.CreatingError, channel)
	}

	friendship := &models.Friendship{}

	friendship.DmChannelID = channel.(models.Channel).ID
	friendship.InitiatorID = models.BigInt(initiatorIDInt)
	friendship.FriendID = models.BigInt(friendIDInt)
	friendship.ID = int64(core.GenerateID())

	return friendshipRepo.CreateFriendship(db, friendship)
}

func (fs *FriendshipService) GetFriendshipRequests(db *gorm.DB, data interface{}) (interface{}, error) {
	userID := data.(string)

	if userID == "" {
		return nil, fs.GenError(fs.InvalidData, nil)
	}

	friendshipRepo := repositories.NewFriendshipRepository()

	return friendshipRepo.GetFriendshipRequests(db, userID)
}

func (fs *FriendshipService) GetFriendships(db *gorm.DB, data interface{}) (interface{}, error) {
	userID := data.(string)

	friendshipRepo := repositories.NewFriendshipRepository()

	return friendshipRepo.GetFriendships(db, userID)
}

func (fs *FriendshipService) DeleteFriendshipRequest(db *gorm.DB, data interface{}) (interface{}, error) {
	id := data.(string)

	friendshipRepo := repositories.NewFriendshipRepository()

	return friendshipRepo.DeleteFriendshipRequest(db, id)
}

func (fs *FriendshipService) DeleteFriendship(db *gorm.DB, data interface{}) (interface{}, error) {
	userID := data.(string)

	friendshipRepo := repositories.NewFriendshipRepository()

	friendship, err := friendshipRepo.GetFriendship(db, userID)
	if err != nil {
		return nil, fs.GenError(fs.GettingError, friendship)
	}

	channelID := strconv.Itoa(int(friendship.(models.Friendship).DmChannelID))

	_, err2 := friendshipRepo.DeleteFriendship(db, friendship)
	if err2 != nil {
		return nil, fs.GenError(fs.DeleteError, nil)
	}

	channelService := NewChannelService()
	_, err3 := channelService.DeleteChannel(db, channelID)

	if err3 != nil {
		return nil, fs.GenError(fs.DeleteError, nil)
	}

	return nil, nil
}

func NewFriendshipService() *FriendshipService {
	friendshipService := &FriendshipService{}

	return &FriendshipService{
		Service: *core.NewService(
			[]core.ServiceMethod{
				friendshipService.CreateFriendship,
				friendshipService.GetFriendships,
				friendshipService.DeleteFriendship,
				friendshipService.CreateFriendshipRequest,
				friendshipService.GetFriendshipRequests,
				friendshipService.DeleteFriendshipRequest,
			},
		),
	}
}
