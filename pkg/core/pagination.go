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

func NewPagination(pageSize, pageNumber int) *Pagination {
	return &Pagination{
		PageSize: pageSize,
		PageNumber: pageNumber,
	}
}
