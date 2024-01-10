package core

import (
	"go.uber.org/zap"

	"gorm.io/gorm"
)

type RepositoryMethod func(*gorm.DB, interface{}) (interface{}, error)

type Repository struct {
	logger zap.Logger
	target interface{}
	methods []RepositoryMethod
}

func NewRepository(model interface{}, methods []RepositoryMethod) *Repository {
	return &Repository{
		target: model,
		methods: methods,
	}
}
