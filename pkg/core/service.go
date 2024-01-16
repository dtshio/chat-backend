package core

import (
	"gorm.io/gorm"
)

type ServiceMethod func(*gorm.DB, interface{}) (interface{}, error)

type Service struct {
	methods []ServiceMethod

	GenError func(string, interface{}) error
	InvalidData string
	CreatingError string
	GettingError string
	DuplicateError string
	UpdateError string
	DeleteError string
}

func NewService(methods []ServiceMethod) *Service {
	return &Service{
		methods: methods,

		GenError: GenError,
		InvalidData: InvalidData,
		CreatingError: CreatingError,
		GettingError: GettingError,
		DuplicateError: DuplicateError,
		UpdateError: UpdateError,
		DeleteError: DeleteError,
	}
}
