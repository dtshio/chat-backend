package models

type Channel struct {
	ID uint64 `json:"id" gorm:"primaryKey"`
	Type string `json:"type" gorm:"not null"`
}
