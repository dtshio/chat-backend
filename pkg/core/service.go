package core

type ServiceMethod func(interface{}) (interface{}, error)

type Service struct {
	methods []ServiceMethod
}

func NewService(methods []ServiceMethod) *Service {
	return &Service{
		methods,
	}
}
