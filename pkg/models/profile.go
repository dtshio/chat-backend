package models

type Profile struct {
  ID BigInt `json:"id" gorm:"primaryKey"`
  UserID BigInt `json:"user_id" gorm:"not null REFERENCES users(id)"`
  Username string `json:"username" gorm:"not null"`
}
