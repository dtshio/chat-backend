package core

import (
	"os"
	"time"

	"github.com/disgoorg/snowflake/v2"
	"github.com/pierrec/xxHash/xxHash32"
)

func GenerateID() uint64 {
	return uint64(snowflake.New(time.Now()))
}

var secretKey = []byte(os.Getenv("CHAT_SECRET_KEY"))

func HashID(firstID string, secondID string) uint32 {
	var id string

	if firstID > secondID {
		id = firstID + secondID
	} else {
		id = secondID + firstID
	}

	var seed uint32
	for _, byte := range secretKey {
		seed += uint32(byte)
	}

	hasher := xxHash32.New(seed)
	hasher.Write([]byte(id))
	hash := hasher.Sum32()
	hasher.Reset()

	return hash
}
