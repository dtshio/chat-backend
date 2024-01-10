package core

import (
	"net/http"
)

type ControllerMethod http.HandlerFunc

type Controller struct {
	methods []ControllerMethod
}

func NewController(methods []ControllerMethod) *Controller {
	return &Controller{
		methods: methods,
	}
}
