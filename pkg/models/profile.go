package models

type Profile struct {
  ID BigInt `json:"id" gorm:"primaryKey"`
  UserID BigInt `json:"user_id" gorm:"not null REFERENCES users(id)"`
  Username string `json:"username" gorm:"not null"`
}


func (p *Profile) BeforeCreateRecord(id BigInt) error {
	if p.ID < 0 {
		p.ID = -p.ID
	}

	return nil
}
