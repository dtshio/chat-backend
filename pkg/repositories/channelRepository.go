package repositories

import (
	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"gorm.io/gorm"
)

type ChannelRepository struct {
	core.Repository
	db *gorm.DB
}

func (cr *ChannelRepository) CreateChannel(data interface {}) (interface {}, error) {
	channel, ok := data.(*models.Channel)

    if !ok {
        return nil, cr.GenError(cr.InvalidData, channel)
    }

	err := channel.BeforeCreateRecord()
	if err != nil {
		return nil, cr.GenError(cr.InvalidData, channel)
	}

	err = cr.db.Table("channels").Create(channel).Error
	if err != nil {
		return nil, cr.GenError(cr.CreateError, channel)
	}

	return *channel, nil
}

func (cr *ChannelRepository) GetChannels(data interface {}) (interface {}, error) {
	channels, ok := data.(*[]models.Channel)

	if !ok {
		return nil, cr.GenError(cr.InvalidData, channels)
	}

	err := cr.db.Table("channels").Find(channels).Error
	if err != nil {
		return nil, cr.GenError(cr.NotFoundError, channels)
	}

	return *channels, err
}

func (cr *ChannelRepository) DeleteChannel(data interface {}) (interface {}, error) {
	id := data.(string)

	err := cr.db.Table("channels").Where("id = ?", id).Delete(&models.Channel{}).Error
	if err != nil {
		return nil, cr.GenError(cr.DeleteError, nil)
	}

	return nil, err
}

func NewChannelRepository(db *gorm.DB) *ChannelRepository {
	return &ChannelRepository{
		Repository: *core.NewRepository(
			&models.Channel{},
			[]core.RepositoryMethod{
				(&ChannelRepository{}).CreateChannel,
				(&ChannelRepository{}).GetChannels,
				(&ChannelRepository{}).DeleteChannel,
			},
		),
		db: db,
	}
}
