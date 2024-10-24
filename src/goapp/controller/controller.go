package controller

import (
	"main/config"
	cApplicationModule "main/controller/app-module"
	cItem "main/controller/item"
	cUser "main/controller/user"
	"main/service"
)

type Controller struct {
	Item              cItem.ItemController
	ItemPage          cItem.ItemPageController
	ApplicationModule cApplicationModule.ApplicationModuleController
	User              cUser.UserController
}

type ControllerOptionFunc func(*Controller)

func NewController(opts ...ControllerOptionFunc) *Controller {
	controller := &Controller{}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

func NewItemController(svc *service.Service) ControllerOptionFunc {
	return func(c *Controller) {
		c.Item = cItem.NewItemController(svc)
	}
}

func NewApplicationModuleController(svc *service.Service) ControllerOptionFunc {
	return func(c *Controller) {
		c.ApplicationModule = cApplicationModule.NewApplicationModuleController(svc)
	}
}

func NewUserController(svc *service.Service) ControllerOptionFunc {
	return func(c *Controller) {
		c.User = cUser.NewUserController(svc)
	}
}

// PAGES

func NewItemPageController(svc *service.Service, conf config.ConfigManager) ControllerOptionFunc {
	return func(c *Controller) {
		c.ItemPage = cItem.NewItemPageController(svc, conf)
	}
}
