package models

type Message struct {
	ID BigInt `json:"id" gorm:"primaryKey"`
	ChannelID BigInt `json:"channel_id" gorm:"not null REFERENCES channels(id)"`
	AuthorID BigInt `json:"author_id" gorm:"not null REFERENCES users(id)"`
	Content string `json:"content" gorm:"not null"`
}
