package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/datsfilipe/pkg/application/redis"
	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"github.com/datsfilipe/pkg/repositories"
	"gorm.io/gorm"
)

type MessageService struct {
	core.Service
}

func (ms *MessageService) CreateMessage(db *gorm.DB, data interface{}) (interface{}, error) {
	message := data.(*models.Message)
	if message == nil {
		return nil, ms.GenError(ms.InvalidData, message)
	}

	messageRepo := repositories.NewMessageRepository()

	newMessage, err := messageRepo.CreateMessage(db, message)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	redis := redis.Open()
	redisMessage, err := json.Marshal(map[string]any{"content": message.Content})
	statusCmd := redis.Publish(ctx, "channel:" + fmt.Sprint(message.ChannelID), redisMessage)

	if statusCmd.Err() != nil {
		return nil, statusCmd.Err()
	}

	return newMessage, nil
}

func (ms *MessageService) GetMessages(db *gorm.DB, data interface{}) (interface{}, error) {
	paylaod := data.(core.Map)
	page := int(paylaod["page"].(float64))
	key := paylaod["channel_id"].(string)

	if key == "" {
		return nil, ms.GenError(ms.InvalidData, key)
	}

	if page < 1 {
		return nil, ms.GenError(ms.InvalidData, page)
	}

	pagination := core.NewPagination(db, 20, page - 1)
	pagination.Key = key

	messageRepo := repositories.NewMessageRepository()
	messages, err := messageRepo.GetMessages(db, pagination)

	return messages, err
}

func NewMessageService() *MessageService {
	messageService := &MessageService{}

	return &MessageService{
		Service: *core.NewService(
			[]core.ServiceMethod{
				messageService.CreateMessage,
				messageService.GetMessages,
			},
		),
	}
}
