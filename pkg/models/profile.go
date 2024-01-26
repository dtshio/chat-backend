package models

import "github.com/datsfilipe/pkg/core"

type Profile struct {
  ID BigInt `json:"id" gorm:"primaryKey"`
  UserID BigInt `json:"user_id" gorm:"not null REFERENCES users(id)"`
  Username string `json:"username" gorm:"not null"`
}

func createID(id BigInt) BigInt {
	if id != 0 {
		return id
	}

	return BigInt(core.GenerateID())
}

func (p *Profile) BeforeCreateRecord(id BigInt) error {
	p.ID = createID(id)

	if p.ID < 0 {
		p.ID = -p.ID
	}

	return nil
}
