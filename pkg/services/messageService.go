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
)

type MessageService struct {
	core.Service
	log *zap.Logger
	repo *repositories.MessageRepository
	service *UserService
}

func (ms *MessageService) CreateMessage(data interface{}) (interface{}, error) {
	message := data.(*models.Message)
	if message == nil {
		return nil, ms.GenError(ms.InvalidData, message)
	}

	messageRecord, err := ms.repo.CreateMessage(message)
	if err != nil {
		ms.log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	authorIDStr := fmt.Sprint(message.AuthorID)

	authorRecords, err := ms.service.GetProfile(authorIDStr)
	if err != nil {
		ms.log.Warn("Warn", zap.Any("Warn", err.Error()))
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
		ms.log.Warn("Warn", zap.Any("Warn", err.Error()))
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

func (ms *MessageService) GetMessages(data interface{}) (interface{}, error) {
	paylaod := data.(core.Map)

	page := int(paylaod["page"].(float64))
	key := paylaod["channel_id"].(string)

	if key == "" {
		return nil, ms.GenError(ms.InvalidData, key)
	}

	if page < 1 {
		return nil, ms.GenError(ms.InvalidData, page)
	}

	pagination := core.NewPagination(20, page - 1)
	pagination.Key = key

	messageRecords, err := ms.repo.GetMessages(pagination)
	if err != nil {
		ms.log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	authorRecords := make([]models.Profile, 0)

	for _, message := range messageRecords.([]models.Message) {
		authorIDStr := fmt.Sprint(message.AuthorID)

		authorRecord, err := ms.service.GetProfile(authorIDStr)
		if err != nil {
			ms.log.Warn("Warn", zap.Any("Warn", err.Error()))
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

func NewMessageService(log *zap.Logger, repo *repositories.MessageRepository, userService *UserService) *MessageService {
	return &MessageService{
		log: log,
		repo: repo,
		service: userService,
	}
}
