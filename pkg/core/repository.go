package core

type RepositoryMethod func(interface{}) (interface{}, error)

type Repository struct {
	target interface{}
	methods []RepositoryMethod

	GenError func(string, interface{}) error
	InvalidData string
	CreateError string
	NotFoundError string
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
		CreateError: CreateError,
		NotFoundError: NotFoundError,
		DuplicateError: DuplicateError,
		UpdateError: UpdateError,
		DeleteError: DeleteError,
	}
}
