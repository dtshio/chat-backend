package models

type Profile struct {
  ID string `json:"id" gorm:"primaryKey"`
  UserID string `json:"user_id" gorm:"not null REFERENCES users(id)"`
  Username string `json:"username" gorm:"not null"`
}
