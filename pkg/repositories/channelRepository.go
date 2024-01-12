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

func (ur *ChannelRepository) CreateChannel(db *gorm.DB, data interface {}) (interface {}, error) {
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

func NewChannelRepository() *ChannelRepository {
	channelRepo := &ChannelRepository{}
	return &ChannelRepository{
		Repository: *core.NewRepository(
			&models.Channel{},
			[]core.RepositoryMethod{
				channelRepo.CreateChannel,
			},
		),
	}
}

