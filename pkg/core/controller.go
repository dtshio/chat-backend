package core

import (
	"net/http"

	"gorm.io/gorm"
)

type ControllerMethod http.HandlerFunc

type Controller struct {
	methods []ControllerMethod
	db *gorm.DB
}

func (c *Controller) SetDB(db *gorm.DB) {
	c.db = db
}

func NewController(methods []ControllerMethod) *Controller {
	return &Controller{
		methods: methods,
	}
}
