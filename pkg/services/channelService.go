package services

import (
	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"github.com/datsfilipe/pkg/repositories"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ChannelService struct {
	core.Service
}

func (cs *ChannelService) CreateChannel(db *gorm.DB, log *zap.Logger, data interface{}) (interface{}, error) {
	channel := data.(*models.Channel)

	channel.Type = "DIRECT_MESSAGE"

	repo := repositories.NewChannelRepository()

	dbRecord, err := repo.CreateChannel(db, channel)
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecord, nil
}

func (cs *ChannelService) GetChannels(db *gorm.DB, log *zap.Logger, data interface{}) (interface{}, error) {
	channels := data.(*[]models.Channel)
	if channels == nil {
		return nil, cs.GenError(cs.InvalidData, channels)
	}

	repo := repositories.NewChannelRepository()

	dbRecords, err := repo.GetChannels(db, channels)
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecords, nil
}

func (cs *ChannelService) DeleteChannel(db *gorm.DB, log *zap.Logger, data interface{}) (interface{}, error) {
	id, ok := data.(string)

	if !ok || id == "" {
		return nil, cs.GenError(cs.InvalidData, nil)
	}

	repo := repositories.NewChannelRepository()

	dbRecord, err := repo.DeleteChannel(db, id)
	if err != nil {
		log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecord, nil
}

func NewChannelService() *ChannelService {
	service := &ChannelService{}

	return &ChannelService{
		Service: *core.NewService(
			[]core.ServiceMethod{
				service.CreateChannel,
				service.GetChannels,
				service.DeleteChannel,
			},
		),
	}
}

