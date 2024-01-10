package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strings"
)

func SignToken(payload string) string {
	key := []byte("3735928559")

	h := hmac.New(sha256.New, key)
	h.Write([]byte(payload))

	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

func VerifyToken(token string) interface{} {
	payload, signature := strings.Split(token, ".")[0], strings.Split(token, ".")[1]

	if SignToken(payload) == signature {
		return payload
	}

	return false
}
