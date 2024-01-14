package core

import (
	"fmt"
	"strconv"
	"time"

	"github.com/disgoorg/snowflake/v2"
)

func GenerateID() uint64 {
	return uint64(snowflake.New(time.Now()))
}

func StringToUint64(str string) (uint64, error) {
	val, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to convert string to uint64: %v", err)
	}
	return val, nil
}
