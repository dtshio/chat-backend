package models

type Users struct {
	ID string `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"not null"`
	Password string `json:"password" gorm:"not null"`
	Email string `json:"email" gorm:"not null"`
}
