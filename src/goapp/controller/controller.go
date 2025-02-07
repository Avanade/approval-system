package controller

import (
	"main/config"
	cApplicationModule "main/controller/app-module"
	cAuthentication "main/controller/authentication"
	cFallback "main/controller/fallback"
	cInvolvement "main/controller/involvement"
	cIPDiscloure "main/controller/ip-disclosure"
	cItem "main/controller/item"
	cItemActivity "main/controller/item-activity"
	cUser "main/controller/user"
	"main/service"
)

type Controller struct {
	ApplicationModule  cApplicationModule.ApplicationModuleController
	AuthenticationPage cAuthentication.AuthenticationPageController
	Fallback           cFallback.FallbackController
	Involvement        cInvolvement.InvolvementController
	IPDisclosure       cIPDiscloure.IpDisclosureController
	IPDisclourePage    cIPDiscloure.IpDisclosurePageController
	Item               cItem.ItemController
	ItemActivity       cItemActivity.ItemActivityController
	ItemPage           cItem.ItemPageController
	User               cUser.UserController
}

type ControllerOptionFunc func(*Controller)

func NewController(opts ...ControllerOptionFunc) *Controller {
	controller := &Controller{}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

func NewApplicationModuleController(svc *service.Service) ControllerOptionFunc {
	return func(c *Controller) {
		c.ApplicationModule = cApplicationModule.NewApplicationModuleController(svc)
	}
}

func NewAuthenticationController(svc *service.Service) ControllerOptionFunc {
	return func(c *Controller) {
		c.AuthenticationPage = cAuthentication.NewAuthenticationController(svc)
	}
}

func NewInvolvementController(svc *service.Service) ControllerOptionFunc {
	return func(c *Controller) {
		c.Involvement = cInvolvement.NewInvolvementController(svc)
	}
}

func NewIPDisclosureController(svc *service.Service, conf config.ConfigManager) ControllerOptionFunc {
	return func(c *Controller) {
		c.IPDisclosure = cIPDiscloure.NewIpDisclosureController(svc, conf)
	}
}

func NewItemController(svc *service.Service) ControllerOptionFunc {
	return func(c *Controller) {
		c.Item = cItem.NewItemController(svc)
	}
}

func NewItemActivityController(svc *service.Service, conf config.ConfigManager) ControllerOptionFunc {
	return func(c *Controller) {
		c.ItemActivity = cItemActivity.NewItemActivityController(svc, conf)
	}
}

func NewUserController(svc *service.Service) ControllerOptionFunc {
	return func(c *Controller) {
		c.User = cUser.NewUserController(svc)
	}
}

// PAGES

func NewIPDisclosurePageController(svc *service.Service) ControllerOptionFunc {
	return func(c *Controller) {
		c.IPDisclourePage = cIPDiscloure.NewIpDisclosurePageController(svc)
	}
}

func NewItemPageController(svc *service.Service, conf config.ConfigManager) ControllerOptionFunc {
	return func(c *Controller) {
		c.ItemPage = cItem.NewItemPageController(svc, conf)
	}
}

func NewFallbackController(svc *service.Service) ControllerOptionFunc {
	return func(c *Controller) {
		c.Fallback = cFallback.NewFallbackController(svc)
	}
}
