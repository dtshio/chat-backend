package repositories

import (
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
        return nil, cr.GenError(cr.InvalidData, channel)
    }

	err := db.Table("channels").Create(channel).Error
	if err != nil {
		return nil, cr.GenError(cr.CreatingError, channel)
	}

	return *channel, err
}

func (cr *ChannelRepository) GetChannels(db *gorm.DB, data interface {}) (interface {}, error) {
	channels, ok := data.(*[]models.Channel)
	if !ok {
		return nil, cr.GenError(cr.InvalidData, channels)
	}

	err := db.Table("channels").Find(channels).Error
	if err != nil {
		return nil, cr.GenError(cr.GettingError, channels)
	}

	return *channels, err
}

func (cr *ChannelRepository) DeleteChannel(db *gorm.DB, data interface {}) (interface {}, error) {
	id := data.(string)

	err := db.Table("channels").Where("id = ?", id).Delete(&models.Channel{}).Error
	if err != nil {
		return nil, cr.GenError(cr.DeleteError, nil)
	}

	return nil, err
}

func NewChannelRepository() *ChannelRepository {
	channelRepo := &ChannelRepository{}
	return &ChannelRepository{
		Repository: *core.NewRepository(
			&models.Channel{},
			[]core.RepositoryMethod{
				channelRepo.CreateChannel,
				channelRepo.GetChannels,
				channelRepo.DeleteChannel,
			},
		),
	}
}

