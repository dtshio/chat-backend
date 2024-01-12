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

func NewMessageRepository() *MessageRepository {
	MessageRepo := &MessageRepository{}
	return &MessageRepository{
		Repository: *core.NewRepository(
			&models.Message{},
			[]core.RepositoryMethod{
				MessageRepo.CreateMessage,
			},
		),
	}
}
