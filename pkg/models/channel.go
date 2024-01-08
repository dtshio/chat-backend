package models

type Channel struct {
	ID string `json:"id" gorm:"primaryKey"`
	Type string `json:"type" gorm:"not null"`
}
