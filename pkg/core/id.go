package core

import (
	"os"
	"time"

	"github.com/cespare/xxhash/v2"
	"github.com/disgoorg/snowflake/v2"
)

func GenerateID() int64 {
	return int64(snowflake.New(time.Now()))
}

var secretKey = []byte(os.Getenv("CHAT_SECRET_KEY"))

func HashID(firstID string, secondID string) int64 {
	var id string

	if firstID > secondID {
		id = firstID + secondID
	} else {
		id = secondID + firstID
	}

	hasher := xxhash.New()
	hasher.Write([]byte(id))
	hashID :=  hasher.Sum64()

	return int64(hashID)
}
