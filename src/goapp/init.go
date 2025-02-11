package main

import (
	"main/config"
	"main/infrastructure/database"
	"main/infrastructure/session"
	"main/middleware"
	"main/router"

	c "main/controller"
	r "main/repository"
	s "main/service"
	t "main/timed-jobs"
)

var (
	conf config.ConfigManager   = config.NewEnvConfigManager()
	db   database.Database      = database.NewDatabase(conf)
	cs   session.ConnectSession = session.NewSession(conf)

	repo = r.NewRepository(
		r.NewApplication(&db),
		r.NewApplicationModule(&db),
		r.NewApprovalRequestApprover(&db),
		r.NewInvolvement(&db),
		r.NewIPDRequest(&db),
		r.NewIpdrInvolvement(&db),
		r.NewItem(&db),
		r.NewItemActivity(&db),
	)

	svc = s.NewService(
		s.NewApplicationModuleService(repo),
		s.NewApplicationService(repo),
		s.NewApprovalRequestApproverService(repo),
		s.NewAuthenticatorService(conf, &cs),
		s.NewEmailService(conf),
		s.NewInvolvementService(repo),
		s.NewIPDisclosureRequestService(repo),
		s.NewItemService(repo, conf),
		s.NewItemActivityService(repo),
		s.NewMsGraphService(conf),
		s.NewTemplateService(conf),
	)

	ctrl = c.NewController(
		c.NewApplicationModuleController(svc),
		c.NewAuthenticationController(svc),
		c.NewFallbackController(svc),
		c.NewInvolvementController(svc),
		c.NewIPDisclosureController(svc, conf),
		c.NewIPDisclosurePageController(svc),
		c.NewItemActivityController(svc, conf),
		c.NewItemController(svc),
		c.NewItemPageController(svc, conf),
		c.NewUserController(svc),
	)

	timedJobs = t.NewTimedJobs(svc, conf)

	m          middleware.Middleware = middleware.NewMiddleware(svc)
	httpRouter router.Router         = router.NewMuxRouter(ctrl, conf, &m)
)
