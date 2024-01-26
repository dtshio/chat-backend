package models

import "github.com/datsfilipe/pkg/core"

type Message struct {
	ID BigInt `json:"id" gorm:"primaryKey"`
	ChannelID BigInt `json:"channel_id" gorm:"not null REFERENCES channels(id)"`
	AuthorID any `json:"author_id" gorm:"REFERENCES profiles(id)"`
	Content string `json:"content" gorm:"not null"`
}

func (m *Message) BeforeCreateRecord() error {
	m.ID = BigInt(core.GenerateID())

	if m.ID < 0 {
		m.ID = -m.ID
	}

	return nil
}
