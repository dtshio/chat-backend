package core

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ServiceMethod func(*gorm.DB, *zap.Logger, interface{}) (interface{}, error)

type Service struct {
	methods []ServiceMethod

	GenError func(string, interface{}) error
	InvalidData string
	CreateError string
	NotFoundError string
	DuplicateError string
	UpdateError string
	DeleteError string
}

func NewService(methods []ServiceMethod) *Service {
	return &Service{
		methods: methods,

		GenError: GenError,
		InvalidData: InvalidData,
		CreateError: CreateError,
		NotFoundError: NotFoundError,
		DuplicateError: DuplicateError,
		UpdateError: UpdateError,
		DeleteError: DeleteError,
	}
}
