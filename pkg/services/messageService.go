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
		return nil, fmt.Errorf("Data is not a valid Message")
	}

	messageRepo := repositories.NewMessageRepository()

	message.ID = core.GenerateID()

	newMessage, err := messageRepo.CreateMessage(db, message)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	redis := redis.Open()
	redisMessage, err := json.Marshal(fmt.Sprint("{\"content\": \"", message.Content, "\"}"))
	statusCmd := redis.Set(ctx, fmt.Sprint(message.ChannelID), redisMessage, 0)

	if statusCmd.Err() != nil {
		return nil, statusCmd.Err()
	}

	return newMessage, nil
}

func (ms *MessageService) GetMessages(db *gorm.DB, data interface{}) (interface{}, error) {
	paylaod := data.(*core.GetMessagesPayload)
	if paylaod == nil {
		return nil, fmt.Errorf("Data is not a valid GetMessagesPayload")
	}

	if paylaod.Page < 1 {
		return nil, fmt.Errorf("Page number must be greater than 0")
	}

	messageRepo := repositories.NewMessageRepository()
	pagination := core.NewPagination(db, 20, paylaod.Page - 1)
	pagination.Key = paylaod.ChannelID

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
