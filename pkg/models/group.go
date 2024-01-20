package models

import "github.com/datsfilipe/pkg/core"

type Group struct {
	ID BigInt `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"not null"`
	OwnerID BigInt `json:"owner_id" gorm:"not null REFERENCES users(id)"`
	ChannelID BigInt `json:"channel_id" gorm:"not null REFERENCES channels(id)"`
}

type GroupMember struct {
	ID BigInt `json:"id" gorm:"primaryKey"`
	GroupID BigInt `json:"group_id" gorm:"not null REFERENCES groups(id)"`
	ProfileID BigInt `json:"profile_id" gorm:"not null REFERENCES profiles(id)"`
	ChannelID BigInt `json:"channel_id" gorm:"not null REFERENCES channels(id)"`
}

func (g *Group) BeforeCreateRecord() error {
	g.ID = BigInt(core.GenerateID())

	if g.ID < 0 {
		g.ID = -g.ID
	}

	return nil
}

func (gm *GroupMember) BeforeCreateRecord() error {
	gm.ID = BigInt(core.GenerateID())

	if gm.ID < 0 {
		gm.ID = -gm.ID
	}

	return nil
}
