package services

import (
	"fmt"

	"github.com/datsfilipe/pkg/application/database"
	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"github.com/datsfilipe/pkg/repositories"
)

type ChannelService struct {
	core.Service
}

func (us *ChannelService) CreateChannel(data interface{}) (interface{}, error) {
	channel := data.(*models.Channel)
	if channel == nil {
		return nil, fmt.Errorf("Data is not a valid Channel")
	}

	db, err := database.Open()
	if err != nil {
		return nil, err
	}

	channelRepo := repositories.NewChannelRepository()

	channel.ID = core.GenerateID()
	channel.Type = "DIRECT_MESSAGE"

	return channelRepo.CreateChannel(db, channel)
}

func NewChannelService() *ChannelService {
	channelService := &ChannelService{}

	return &ChannelService{
		Service: *core.NewService(
			[]core.ServiceMethod{
				channelService.CreateChannel,
			},
		),
	}
}

