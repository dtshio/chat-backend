package models

import (
	"strconv"

	"github.com/datsfilipe/pkg/core"
)

type FriendshipBase struct {
	ID          BigInt `json:"id" gorm:"primaryKey"`
	InitiatorID BigInt `json:"initiator_id" gorm:"not null REFERENCES profiles(id)"`
	FriendID    BigInt `json:"friend_id" gorm:"not null REFERENCES profiles(id)"`
}

type FriendshipRequest struct {
	FriendshipBase
	Accepted bool `json:"accepted" gorm:"default:false"`
}

type Friendship struct {
	FriendshipBase
	DmChannelID BigInt `json:"dm_channel_id" gorm:"not null REFERENCES channels(id)"`
}

func (fr *FriendshipRequest) BeforeCreateRecord() error {
	initiatorID := strconv.Itoa(int(fr.InitiatorID))
	friendID := strconv.Itoa(int(fr.FriendID))

	id := core.HashID(initiatorID, friendID)

	if id < 0 {
		id = -id
	}

	fr.ID = BigInt(id)

	return nil
}

func (f *Friendship) BeforeCreateRecord() error {
	f.ID = BigInt(core.GenerateID())

	if f.ID < 0 {
		f.ID = -f.ID
	}

	return nil
}
