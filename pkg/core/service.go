package core

import "gorm.io/gorm"

type ServiceMethod func(*gorm.DB, interface{}) (interface{}, error)

type Service struct {
	methods []ServiceMethod
}

func NewService(methods []ServiceMethod) *Service {
	return &Service{
		methods,
	}
}
