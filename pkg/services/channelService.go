package services

import (
	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"github.com/datsfilipe/pkg/repositories"
	"go.uber.org/zap"
)

type ChannelService struct {
	core.Service
	log *zap.Logger
	repo *repositories.ChannelRepository
}

func (cs *ChannelService) CreateChannel(data interface{}) (interface{}, error) {
	channel := data.(*models.Channel)

	dbRecord, err := cs.repo.CreateChannel(channel)
	if err != nil {
		cs.log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecord, nil
}

func (cs *ChannelService) GetChannels(data interface{}) (interface{}, error) {
	channels := data.(*[]models.Channel)
	if channels == nil {
		return nil, cs.GenError(cs.InvalidData, channels)
	}

	dbRecords, err := cs.repo.GetChannels(channels)
	if err != nil {
		cs.log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecords, nil
}

func (cs *ChannelService) DeleteChannel(data interface{}) (interface{}, error) {
	id, ok := data.(string)

	if !ok || id == "" {
		return nil, cs.GenError(cs.InvalidData, nil)
	}

	dbRecord, err := cs.repo.DeleteChannel(id)
	if err != nil {
		cs.log.Warn("Warn", zap.Any("Warn", err.Error()))
		return nil, err
	}

	return dbRecord, nil
}

func NewChannelService(log *zap.Logger, repo *repositories.ChannelRepository) *ChannelService {
	return &ChannelService{
		Service: *core.NewService(),
		log: log,
		repo: repo,
	}
}
