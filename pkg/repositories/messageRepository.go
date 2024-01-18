package repositories

import (
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
        return nil, mr.GenError(mr.InvalidData, message)
    }

	err := message.BeforeCreateRecord()
	if err != nil {
		return nil, mr.GenError(mr.InvalidData, message)
	}

	err1 := db.Table("messages").Create(message).Error
	if err1 != nil {
		return nil, mr.GenError(mr.CreatingError, message)
	}

	return *message, nil
}

func (mr *MessageRepository) GetMessages(db *gorm.DB, data interface{}) (interface{}, error) {
    pagination, ok := data.(*core.Pagination)
    if !ok {
        return nil, mr.GenError(mr.InvalidData, pagination)
    }

    var messages []models.Message
    offset := pagination.PageSize * (pagination.PageNumber)
    
    db = db.Table("messages").Order("created_at DESC").Offset(offset).Limit(pagination.PageSize)
    db = db.Where("channel_id = ?", pagination.Key)
    
    err := db.Find(&messages).Error
    if err != nil {
        return nil, mr.GenError(mr.GettingError, messages)
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
