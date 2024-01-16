package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"os"
	"strings"
)

var secretKey = []byte(os.Getenv("CHAT_SECRET_KEY"))

func SignToken(payload string) string {
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(payload))

	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

func VerifyToken(token string) bool {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return false
	}

	payload, signature := parts[0], parts[1]

	if SignToken(payload) == signature {
		return true
	}

	return false
}
