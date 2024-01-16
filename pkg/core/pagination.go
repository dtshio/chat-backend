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

func NewPagination(db *gorm.DB, pageSize, pageNumber int) *Pagination {
	return &Pagination{
		DB: db,
		PageSize: pageSize,
		PageNumber: pageNumber,
	}
}
