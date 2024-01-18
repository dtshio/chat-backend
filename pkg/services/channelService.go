package services

import (
	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"github.com/datsfilipe/pkg/repositories"
	"gorm.io/gorm"
)

type ChannelService struct {
	core.Service
}

func (cs *ChannelService) CreateChannel(db *gorm.DB, data interface{}) (interface{}, error) {
	channel := data.(*models.Channel)

	channel.ID = models.BigInt(core.GenerateID())
	channel.Type = "DIRECT_MESSAGE"

	channelRepo := repositories.NewChannelRepository()
	return channelRepo.CreateChannel(db, channel)
}

func (cs *ChannelService) GetChannels(db *gorm.DB, data interface{}) (interface{}, error) {
	channels := data.(*[]models.Channel)
	if channels == nil {
		return nil, cs.GenError(cs.InvalidData, channels)
	}

	channelRepo := repositories.NewChannelRepository()

	return channelRepo.GetChannels(db, channels)
}

func (cs *ChannelService) DeleteChannel(db *gorm.DB, data interface{}) (interface{}, error) {
	id, ok := data.(string)

	if !ok || id == "" {
		return nil, cs.GenError(cs.InvalidData, nil)
	}

	channelRepo := repositories.NewChannelRepository()

	return channelRepo.DeleteChannel(db, id)
}

func NewChannelService() *ChannelService {
	channelService := &ChannelService{}

	return &ChannelService{
		Service: *core.NewService(
			[]core.ServiceMethod{
				channelService.CreateChannel,
				channelService.GetChannels,
			},
		),
	}
}

