package models

type Message struct {
	ID string `json:"id" gorm:"primaryKey"`
	ChannelID string `json:"channel_id" gorm:"not null REFERENCES channels(id)"`
	AuthorID string `json:"author_id" gorm:"not null REFERENCES users(id)"`
	Content string `json:"content" gorm:"not null"`
}
