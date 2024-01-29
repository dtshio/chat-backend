package repositories

import (
	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"gorm.io/gorm"
)

type FriendshipRepository struct {
	core.Repository
	db *gorm.DB
}

func (fr *FriendshipRepository) CreateFriendshipRequest(data interface {}) (interface {}, error) {
	friendshipRequest, ok := data.(*models.FriendshipRequest)
    if !ok {
        return nil, fr.GenError(fr.InvalidData, friendshipRequest)
    }

	err := friendshipRequest.BeforeCreateRecord()
	if err != nil {
		return nil, fr.GenError(fr.InvalidData, friendshipRequest)
	}

	err = fr.db.Table("friendship_requests").Create(friendshipRequest).Error
	if err != nil {
		return nil, fr.GenError(fr.CreateError, friendshipRequest)
	}

	return *friendshipRequest, nil
}

func (fr *FriendshipRepository) CreateFriendship(data interface {}) (interface {}, error) {
	friendship, ok := data.(*models.Friendship)

	if !ok {
		return nil, fr.GenError(fr.InvalidData, friendship)
	}

	err := friendship.BeforeCreateRecord()
	if err != nil {
		return nil, fr.GenError(fr.InvalidData, friendship)
	}

	err = fr.db.Table("friendships").Create(friendship).Error
	if err != nil {
		return nil, fr.GenError(fr.CreateError, friendship)
	}

	return *friendship, nil
}

func (fr *FriendshipRepository) GetFriendships(data interface {}) (interface {}, error) {
	userID := data.(string)

	friendships := &[]models.Friendship{}

	err := fr.db.Table("friendships").Where("initiator_id = ? OR friend_id = ?", userID, userID).Find(&friendships).Error
	if err != nil {
		return nil, err
	}
	
	if len(*friendships) == 0 {
		return nil, nil
	}

	return *friendships, nil
}

func (fr *FriendshipRepository) GetFriendship(data interface {}) (interface {}, error) {
	id, ok := data.(string)

	if !ok {
		return nil, fr.GenError(fr.InvalidData, nil)
	}

	friendship := &models.Friendship{}
	
	err := fr.db.Table("friendships").Where("id = ?", id).First(friendship).Error
	if err != nil {
		return nil, fr.GenError(fr.NotFoundError, friendship)
	}

	return *friendship, err
}

func (fr *FriendshipRepository) DeleteFriendship(data interface {}) (interface {}, error) {
	friendship, ok := data.(models.Friendship)

	if !ok {
		return nil, fr.GenError(fr.InvalidData, nil)
	}

	err := fr.db.Table("friendships").Delete(friendship).Error
	if err != nil {
		return nil, fr.GenError(fr.DeleteError, friendship)
	}

	return nil, err
}
	
func (fr *FriendshipRepository) GetFriendshipRequests(data interface {}) (interface {}, error) {
	userID := data.(string)

	if userID == "" {
		return nil, fr.GenError(fr.InvalidData, nil)
	}

	friendshipRequests := &[]models.FriendshipRequest{}

	err := fr.db.Table("friendship_requests").Where("(initiator_id = ? OR friend_id = ?) AND accepted = false", userID, userID).Find(friendshipRequests).Error
	if err != nil {
		return nil, fr.GenError(fr.NotFoundError, friendshipRequests)
	}

	if len(*friendshipRequests) == 0 {
		return nil, fr.GenError(fr.NotFoundError, friendshipRequests)
	}

	return *friendshipRequests, err
}

func (fr *FriendshipRepository) GetFriendshipRequest(data interface {}) (interface {}, error) {
	id, ok := data.(*models.BigInt)

	if !ok {
		return nil, fr.GenError(fr.InvalidData, id)
	}

	friendshipRequest := &models.FriendshipRequest{}

	err := fr.db.Table("friendship_requests").Where("id = ?", id).First(friendshipRequest).Error
	if err != nil {
		return nil, fr.GenError(fr.NotFoundError, friendshipRequest)
	}

	return *friendshipRequest, err
}

func (fr *FriendshipRepository) DeleteFriendshipRequest(data interface {}) (interface {}, error) {
	requestID, ok := data.(string)

	if !ok {
		return nil, fr.GenError(fr.InvalidData, nil)
	}

	friendshipRequest := &models.FriendshipRequest{}

	err := fr.db.Table("friendship_requests").Where("id = ? AND accepted = false", requestID).Delete(friendshipRequest).Error
	if err != nil {
		return nil, fr.GenError(fr.DeleteError, friendshipRequest)
	}

	return nil, err
}

func (fr *FriendshipRepository) AcceptFriendshipRequest(data interface {}) (interface {}, error) {
	requestID, ok := data.(int64)

	if !ok {
		return nil, fr.GenError(fr.InvalidData, nil)
	}

	friendshipRequest := &models.FriendshipRequest{}

	err := fr.db.Table("friendship_requests").Where("id = ?", requestID).Update("accepted", true).Error
	if err != nil {
		return nil, fr.GenError(fr.DeleteError, friendshipRequest)
	}

	return nil, err
}

func NewFriendshipRepository(db *gorm.DB) *FriendshipRepository {
	return &FriendshipRepository{
		Repository: *core.NewRepository(),
		db: db,
	}
}
