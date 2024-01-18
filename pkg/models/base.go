package models

import (
	"encoding/json"
	"strconv"
)

type BigInt int64

func (integer BigInt) MarshalJSON() ([]byte, error) {
    return json.Marshal(strconv.Itoa(int(integer)))
}
