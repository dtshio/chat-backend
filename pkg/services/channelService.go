package services

import (
	"fmt"

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
	if channel == nil {
		return nil, fmt.Errorf("Data is not a valid Channel")
	}

	channelRepo := repositories.NewChannelRepository()

	channel.ID = models.BigInt(core.GenerateID())
	channel.Type = "DIRECT_MESSAGE"

	return channelRepo.CreateChannel(db, channel)
}

func (cs *ChannelService) GetChannels(db *gorm.DB, data interface{}) (interface{}, error) {
	channels := data.(*[]models.Channel)
	if channels == nil {
		return nil, fmt.Errorf("Data is not a valid Channel")
	}

	channelRepo := repositories.NewChannelRepository()

	return channelRepo.GetChannels(db, channels)
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

