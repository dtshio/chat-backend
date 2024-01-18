package models

import "github.com/datsfilipe/pkg/core"

type Channel struct {
	ID BigInt `json:"id" gorm:"primaryKey"`
	Type string `json:"type" gorm:"not null"`
}

func (c *Channel) BeforeCreateRecord() error {
	c.ID = BigInt(core.GenerateID())

	if c.ID < 0 {
		c.ID = -c.ID
	}

	return nil
}
