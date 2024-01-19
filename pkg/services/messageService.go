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

	dbRecord, err := repo.CreateMessage(db, message)
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	ctx := context.Background()
	redis := redis.Open()
	redisMessage, err := json.Marshal(map[string]any{"content": dbRecord.(models.Message).Content})
	status := redis.Publish(ctx, "channel:" + fmt.Sprint(dbRecord.(models.Message).ChannelID), redisMessage)

	if status.Err() != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, ms.GenError(ms.CreateError, status.Err())
	}

	return dbRecord, nil
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

	dbRecords, err := repo.GetMessages(db, pagination)
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecords, nil
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
