package core

import (
	"os"
	"time"

	"github.com/cespare/xxhash/v2"
	"github.com/disgoorg/snowflake/v2"
)

func GenerateID() uint64 {
	return uint64(snowflake.New(time.Now()))
}

var secretKey = []byte(os.Getenv("CHAT_SECRET_KEY"))

func HashID(firstID string, secondID string) uint64 {
	var id string

	if firstID > secondID {
		id = firstID + secondID
	} else {
		id = secondID + firstID
	}

	hasher := xxhash.New()
	hasher.Write([]byte(id))

	return hasher.Sum64()
}
