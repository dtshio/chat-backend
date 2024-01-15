package models

type User struct {
	ID BigInt `json:"id" gorm:"primaryKey"`
	Password string `json:"password" gorm:"not null"`
	Email string `json:"email" gorm:"not null"`
}
