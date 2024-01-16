package models

import (
	"strconv"

	"github.com/cespare/xxhash/v2"
)

type Friendship struct {
	ID BigInt `json:"id" gorm:"primaryKey"`
	InitiatorID BigInt `json:"initiator_id" gorm:"not null"`
	FriendID BigInt `json:"friend_id" gorm:"not null"`
}

func (f *Friendship) BeforeCreate() {
	var id string

	if f.InitiatorID > f.FriendID {
		id = strconv.Itoa(int(f.FriendID)) + strconv.Itoa(int(f.InitiatorID))
	} else {
		id = strconv.Itoa(int(f.InitiatorID)) + strconv.Itoa(int(f.FriendID))
	}

	f.ID = BigInt(xxhash.Sum64String(id))
}
