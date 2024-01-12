package core

import (
	"time"

	"github.com/disgoorg/snowflake/v2"
)

func GenerateID() uint64 {
	return uint64(snowflake.New(time.Now()))
}
