package models

type Channel struct {
	ID BigInt `json:"id" gorm:"primaryKey"`
	Type string `json:"type" gorm:"not null"`
}
