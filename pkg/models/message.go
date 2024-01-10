package models

type Message struct {
	ID uint64 `json:"id" gorm:"primaryKey"`
	ChannelID uint64 `json:"channel_id" gorm:"not null REFERENCES channels(id)"`
	AuthorID uint64 `json:"author_id" gorm:"not null REFERENCES users(id)"`
	Content string `json:"content" gorm:"not null"`
}
