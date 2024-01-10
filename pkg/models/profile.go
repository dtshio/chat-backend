package models

type Profile struct {
  ID uint64 `json:"id" gorm:"primaryKey"`
  UserID uint64 `json:"user_id" gorm:"not null REFERENCES users(id)"`
  Username string `json:"username" gorm:"not null"`
}
