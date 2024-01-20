package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/datsfilipe/pkg/application/redis"
	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"github.com/datsfilipe/pkg/repositories"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MessageService struct {
	core.Service
}

func (ms *MessageService) CreateMessage(db *gorm.DB, log *zap.Logger, data interface{}) (interface{}, error) {
	message := data.(*models.Message)
	if message == nil {
		return nil, ms.GenError(ms.InvalidData, message)
	}

	repo := repositories.NewMessageRepository()

	messageRecord, err := repo.CreateMessage(db, message)
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	userService := NewUserService()

	authorIDStr := fmt.Sprint(message.AuthorID)

	authorRecords, err := userService.GetProfile(db, log, authorIDStr)
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	if len(authorRecords.([]models.Profile)) > 1 {
		return nil, ms.GenError(ms.DuplicateError, authorRecords)
	}

	author := authorRecords.([]models.Profile)[0]

	ctx := context.Background()
	redis := redis.Open()
	redisMessage, _ := json.Marshal(map[string]any{"content": messageRecord.(models.Message).Content, "username": author.Username})
	status := redis.Publish(ctx, "channel:" + fmt.Sprint(messageRecord.(models.Message).ChannelID), redisMessage)

	if status.Err() != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, ms.GenError(ms.CreateError, status.Err())
	}

	messageObject := map[string]any{
		"id": messageRecord.(models.Message).ID,
		"username": author.Username,
		"content": messageRecord.(models.Message).Content,
		"author_id": messageRecord.(models.Message).AuthorID,
		"channel_id": messageRecord.(models.Message).ChannelID,
	}

	return messageObject, nil
}

func (ms *MessageService) GetMessages(db *gorm.DB, log *zap.Logger, data interface{}) (interface{}, error) {
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

	repo := repositories.NewMessageRepository()

	messageRecords, err := repo.GetMessages(db, pagination)
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	authorRecords := make([]models.Profile, 0)

	userService := NewUserService()

	for _, message := range messageRecords.([]models.Message) {
		authorIDStr := fmt.Sprint(message.AuthorID)

		authorRecord, err := userService.GetProfile(db, log, authorIDStr)
		if err != nil {
			log.Warn("Warn", zap.Any("Warn", err.Error()))
			return nil, err
		}

		if len(authorRecord.([]models.Profile)) > 1 {
			return nil, ms.GenError(ms.DuplicateError, authorRecord)
		}

		authorRecords = append(authorRecords, authorRecord.([]models.Profile)[0])
	}

	messages := make([]map[string]any, 0)

	for i, message := range messageRecords.([]models.Message) {
		messageObject := map[string]any{
			"id": message.ID,
			"username": authorRecords[i].Username,
			"content": message.Content,
			"author_id": message.AuthorID,
			"channel_id": message.ChannelID,
		}
		
		messages = append(messages, messageObject)
	}

	return messages, nil
}

func NewMessageService() *MessageService {
	service := &MessageService{}

	return &MessageService{
		Service: *core.NewService(
			[]core.ServiceMethod{
				service.CreateMessage,
				service.GetMessages,
			},
		),
	}
}
