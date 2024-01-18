package models

import (
	"strconv"

	"github.com/datsfilipe/pkg/core"
)

type FriendshipBase struct {
	ID          int64  `json:"id" gorm:"primaryKey"`
	InitiatorID BigInt `json:"initiator_id" gorm:"not null"`
	FriendID    BigInt `json:"friend_id" gorm:"not null"`
}

type FriendshipRequest struct {
	FriendshipBase
	Accepted bool `json:"accepted" gorm:"default:false"`
}

type Friendship struct {
	FriendshipBase
	DmChannelID BigInt
}

func (fr *FriendshipRequest) BeforeCreateRecord() error {
	initiatorID := strconv.Itoa(int(fr.InitiatorID))
	friendID := strconv.Itoa(int(fr.FriendID))

	fr.ID = int64(core.HashID(initiatorID, friendID))

	return nil
}
