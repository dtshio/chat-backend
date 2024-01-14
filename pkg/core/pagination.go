package core

import (
	"gorm.io/gorm"
)

type Pagination struct {
	DB *gorm.DB
	PageSize int
	PageNumber int
	Key string
}

type GetMessagesPayload struct {
	ChannelID string `json:"channel_id"`
	ID string `json:"id"`
	Page int `json:"page"`
}

func NewPagination(db *gorm.DB, pageSize, pageNumber int) *Pagination {
	return &Pagination{
		DB: db,
		PageSize: pageSize,
		PageNumber: pageNumber,
	}
}

func (p *Pagination) ReversePaginate(model interface{}, orderColumn string) ([]interface{}, error) {
	offset := p.PageNumber * p.PageSize

	var results []interface{}
	err := p.DB.Model(model).Order(orderColumn + " DESC").Offset(offset).Limit(p.PageSize).Find(&results).Error

	return results, err
}
