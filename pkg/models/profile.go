package models

import "github.com/datsfilipe/pkg/core"

type Profile struct {
  ID BigInt `json:"id" gorm:"primaryKey"`
  UserID BigInt `json:"user_id" gorm:"not null REFERENCES users(id)"`
  Username string `json:"username" gorm:"not null"`
}


func (p *Profile) BeforeCreateRecord() error {
	p.ID = BigInt(core.GenerateID())

	if p.ID < 0 {
		p.ID = -p.ID
	}

	return nil
}
