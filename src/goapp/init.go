package main

import (
	"main/config"
	"main/infrastructure/database"
	"main/router"

	c "main/controller"
	r "main/repository"
	s "main/service"
)

var (
	conf config.ConfigManager = config.NewEnvConfigManager()
	db   database.Database    = database.NewDatabase(conf)

	repo = r.NewRepository(
		r.NewApplicationModule(db),
		r.NewItem(db),
		r.NewApprovalRequestApprover(db),
	)

	svc = s.NewService(
		s.NewApplicationModuleService(repo),
		s.NewItemService(repo, conf),
		s.NewEmailService(conf),
		s.NewApprovalRequestApproverService(repo),
	)

	ctrl = c.NewController(
		c.NewItemController(svc),
		c.NewApplicationModuleController(svc),
	)

	httpRouter router.Router = router.NewMuxRouter()
)
