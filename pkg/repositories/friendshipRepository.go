package repositories

import (
	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"gorm.io/gorm"
)

type FriendshipRepository struct {
	core.Repository
}

func (fr *FriendshipRepository) CreateFriendshipRequest(db *gorm.DB, data interface {}) (interface {}, error) {
	friendshipRequest, ok := data.(*models.FriendshipRequest)
    if !ok {
        return nil, fr.GenError(fr.InvalidData, friendshipRequest)
    }

	err := friendshipRequest.BeforeCreateRecord()
	if err != nil {
		return nil, fr.GenError(fr.InvalidData, friendshipRequest)
	}

	err1 := db.Table("friendship_requests").Create(friendshipRequest).Error
	if err1 != nil {
		return nil, fr.GenError(fr.CreateError, friendshipRequest)
	}

	return *friendshipRequest, nil
}

func (fr *FriendshipRepository) CreateFriendship(db *gorm.DB, data interface {}) (interface {}, error) {
	friendship, ok := data.(*models.Friendship)

	if !ok {
		return nil, fr.GenError(fr.InvalidData, friendship)
	}

	err := friendship.BeforeCreateRecord()
	if err != nil {
		return nil, fr.GenError(fr.InvalidData, friendship)
	}

	err1 := db.Table("friendships").Create(friendship).Error
	if err1 != nil {
		return nil, fr.GenError(fr.CreateError, friendship)
	}

	return *friendship, nil
}

func (fr *FriendshipRepository) GetFriendships(db *gorm.DB, data interface {}) (interface {}, error) {
	userID, ok := data.(string)
	if !ok || userID == "" {
        return nil, fr.GenError(fr.InvalidData, nil)
    }

	friendships := &[]models.Friendship{}

	err := db.Table("friendships").Where("initiator_id = ? OR friend_id = ?", userID, userID).Find(friendships).Error
	if err != nil {
		return nil, fr.GenError(fr.NotFoundError, friendships)
	}
	
	if len(*friendships) == 0 {
		return nil, fr.GenError(fr.NotFoundError, friendships)
	}

	return *friendships, err
}

func (fr *FriendshipRepository) GetFriendship(db *gorm.DB, data interface {}) (interface {}, error) {
	id, ok := data.(string)

	if !ok {
		return nil, fr.GenError(fr.InvalidData, nil)
	}

	friendship := &models.Friendship{}
	
	err := db.Table("friendships").Where("id = ?", id).First(friendship).Error
	if err != nil {
		return nil, fr.GenError(fr.NotFoundError, friendship)
	}

	return *friendship, err
}

func (fr *FriendshipRepository) DeleteFriendship(db *gorm.DB, data interface {}) (interface {}, error) {
	friendship, ok := data.(models.Friendship)

	if !ok {
		return nil, fr.GenError(fr.InvalidData, nil)
	}

	err := db.Table("friendships").Delete(friendship).Error
	if err != nil {
		return nil, fr.GenError(fr.DeleteError, friendship)
	}

	return nil, err
}
	
func (fr *FriendshipRepository) GetFriendshipRequests(db *gorm.DB, data interface {}) (interface {}, error) {
	userID := data.(string)

	if userID == "" {
		return nil, fr.GenError(fr.InvalidData, nil)
	}

	friendshipRequests := &[]models.FriendshipRequest{}

	err := db.Table("friendship_requests").Where("(initiator_id = ? OR friend_id = ?) AND accepted = false", userID, userID).Find(friendshipRequests).Error
	if err != nil {
		return nil, fr.GenError(fr.NotFoundError, friendshipRequests)
	}

	if len(*friendshipRequests) == 0 {
		return nil, fr.GenError(fr.NotFoundError, friendshipRequests)
	}

	return *friendshipRequests, err
}

func (fr *FriendshipRepository) GetFriendshipRequest(db *gorm.DB, data interface {}) (interface {}, error) {
	id, ok := data.(*models.BigInt)

	if !ok {
		return nil, fr.GenError(fr.InvalidData, id)
	}

	friendshipRequest := &models.FriendshipRequest{}

	err := db.Table("friendship_requests").Where("id = ?", id).First(friendshipRequest).Error
	if err != nil {
		return nil, fr.GenError(fr.NotFoundError, friendshipRequest)
	}

	return *friendshipRequest, err
}

func (fr *FriendshipRepository) DeleteFriendshipRequest(db *gorm.DB, data interface {}) (interface {}, error) {
	requestID, ok := data.(string)

	if !ok {
		return nil, fr.GenError(fr.InvalidData, nil)
	}

	friendshipRequest := &models.FriendshipRequest{}

	err := db.Table("friendship_requests").Where("id = ? AND accepted = false", requestID).Delete(friendshipRequest).Error
	if err != nil {
		return nil, fr.GenError(fr.DeleteError, friendshipRequest)
	}

	return nil, err
}

func (fr *FriendshipRepository) AcceptFriendshipRequest(db *gorm.DB, data interface {}) (interface {}, error) {
	requestID, ok := data.(int64)

	if !ok {
		return nil, fr.GenError(fr.InvalidData, nil)
	}

	friendshipRequest := &models.FriendshipRequest{}

	err := db.Table("friendship_requests").Where("id = ?", requestID).Update("accepted", true).Error
	if err != nil {
		return nil, fr.GenError(fr.DeleteError, friendshipRequest)
	}

	return nil, err
}

func NewFriendshipRepository() *FriendshipRepository {
	FriendshipRepo := &FriendshipRepository{}
	return &FriendshipRepository{
		Repository: *core.NewRepository(
			&models.Friendship{},
			[]core.RepositoryMethod{
				FriendshipRepo.CreateFriendship,
				FriendshipRepo.GetFriendships,
				FriendshipRepo.CreateFriendshipRequest,
				FriendshipRepo.AcceptFriendshipRequest,
				FriendshipRepo.GetFriendshipRequests,
				FriendshipRepo.GetFriendshipRequest,
			},
		),
	}
}
