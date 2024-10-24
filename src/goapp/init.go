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
		r.NewItem(&db),
		r.NewApprovalRequestApprover(&db),
	)

	svc = s.NewService(
		s.NewApplicationService(repo),
		s.NewApplicationModuleService(repo),
		s.NewItemService(repo, conf),
		s.NewEmailService(conf),
		s.NewApprovalRequestApproverService(repo),
		s.NewMsGraphService(conf),
		s.NewTemplateService(conf),
		s.NewAuthenticatorService(conf, &cs),
	)

	ctrl = c.NewController(
		c.NewItemController(svc),
		c.NewApplicationModuleController(svc),
		c.NewUserController(svc),
		c.NewItemPageController(svc, conf),
		c.NewFallbackController(svc),
		c.NewAuthenticationController(svc),
	)

	timedJobs = t.NewTimedJobs(svc, conf)

	m          middleware.Middleware = middleware.NewMiddleware(svc)
	httpRouter router.Router         = router.NewMuxRouter(ctrl, conf)
)
