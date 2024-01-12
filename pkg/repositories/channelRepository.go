package repositories

import (
	"fmt"

	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"gorm.io/gorm"
)

type ChannelRepository struct {
	core.Repository
}

func (cr *ChannelRepository) CreateChannel(db *gorm.DB, data interface {}) (interface {}, error) {
	channel, ok := data.(*models.Channel)
    if !ok {
        return nil, fmt.Errorf("Data is not a valid Channel")
    }

	err := db.Table("channels").Create(channel).Error
	if err != nil {
		return nil, fmt.Errorf("Error creating Channel: %v", err)
	}

	return *channel, err
}

func (cr *ChannelRepository) GetChannels(db *gorm.DB, data interface {}) (interface {}, error) {
	channels, ok := data.(*[]models.Channel)
	if !ok {
		return nil, fmt.Errorf("Data is not a valid Channel")
	}

	err := db.Table("channels").Find(channels).Error
	if err != nil {
		return nil, fmt.Errorf("Error getting Channels: %v", err)
	}

	return *channels, err
}

func NewChannelRepository() *ChannelRepository {
	channelRepo := &ChannelRepository{}
	return &ChannelRepository{
		Repository: *core.NewRepository(
			&models.Channel{},
			[]core.RepositoryMethod{
				channelRepo.CreateChannel,
				channelRepo.GetChannels,
			},
		),
	}
}

