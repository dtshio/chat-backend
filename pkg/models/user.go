package models

import (
	"golang.org/x/crypto/argon2"
	"encoding/base64"
	"os"
)

var secretKey = []byte(os.Getenv("CHAT_SECRET_KEY"))

type User struct {
	ID BigInt `json:"id" gorm:"primaryKey"`
	Password string `json:"password" gorm:"not null"`
	Email string `json:"email" gorm:"not null"`
}

func (u *User) HashPassword(password string) string {
	hashedPassword := argon2.IDKey([]byte(password), secretKey, 1, 64*1024, 4, 32)
	return base64.StdEncoding.EncodeToString(hashedPassword)
}

func (u *User) VerifyPassword(password string, hashedPassword string) bool {
	hashedInput := u.HashPassword(password)
	return hashedInput == hashedPassword
}

func (u *User) BeforeCreate() error {
	u.Password = u.HashPassword(u.Password)
	return nil
}
