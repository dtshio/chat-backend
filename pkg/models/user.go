package models

import (
	"encoding/base64"
	"os"

	"github.com/datsfilipe/pkg/core"
	"golang.org/x/crypto/argon2"
)

var secretKey = []byte(os.Getenv("CHAT_SECRET_KEY"))

type User struct {
	ID BigInt `json:"id" gorm:"primaryKey"`
	Password string `json:"password" gorm:"not null"`
	Email string `json:"email" gorm:"not null"`
	DefaultProfileID BigInt `json:"default_profile_id" gorm:"not null REFERENCES profiles(id)"`
}

func (u *User) HashPassword(password string) string {
	hashedPassword := argon2.IDKey([]byte(password), secretKey, 1, 64*1024, 4, 32)
	return base64.StdEncoding.EncodeToString(hashedPassword)
}

func (u *User) VerifyPassword(password string, hashedPassword string) bool {
	hashedInput := u.HashPassword(password)
	return hashedInput == hashedPassword
}

func (u *User) BeforeCreateRecord() error {
	u.ID = BigInt(core.GenerateID())

	if u.ID < 0 {
		u.ID = -u.ID
	}

	u.Password = u.HashPassword(u.Password)

	return nil
}
