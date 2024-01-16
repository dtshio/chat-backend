package core

import (
	"gorm.io/gorm"
)

type RepositoryMethod func(*gorm.DB, interface{}) (interface{}, error)

type Repository struct {
	target interface{}
	methods []RepositoryMethod

	GenError func(string, interface{}) error
	InvalidData string
	CreatingError string
	GettingError string
	DuplicateError string
	UpdateError string
	DeleteError string
}

func NewRepository(model interface{}, methods []RepositoryMethod) *Repository {
	return &Repository{
		target: model,
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
