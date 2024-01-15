package repositories

import (
	"fmt"

	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"gorm.io/gorm"
)

type MessageRepository struct {
	core.Repository
}

func (mr *MessageRepository) CreateMessage(db *gorm.DB, data interface {}) (interface {}, error) {
	message, ok := data.(*models.Message)
    if !ok {
        return nil, fmt.Errorf("Data is not a valid Message")
    }

	err := db.Table("messages").Create(message).Error
	if err != nil {
		return nil, fmt.Errorf("Error creating Message: %v", err)
	}

	return *message, err
}

func (mr *MessageRepository) GetMessages(db *gorm.DB, data interface{}) (interface{}, error) {
    pagination, ok := data.(*core.Pagination)
    if !ok {
        return nil, fmt.Errorf("Data is not a valid Pagination")
    }

    var messages []models.Message
    offset := pagination.PageSize * (pagination.PageNumber)
    
    db = db.Table("messages").Order("created_at DESC").Offset(offset).Limit(pagination.PageSize)
    db = db.Where("channel_id = ?", pagination.Key)
    
    err := db.Find(&messages).Error
    if err != nil {
        return nil, fmt.Errorf("Error getting Messages: %v", err)
    }

	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

    return messages, err
}

func NewMessageRepository() *MessageRepository {
	MessageRepo := &MessageRepository{}
	return &MessageRepository{
		Repository: *core.NewRepository(
			&models.Message{},
			[]core.RepositoryMethod{
				MessageRepo.CreateMessage,
				MessageRepo.GetMessages,
			},
		),
	}
}
