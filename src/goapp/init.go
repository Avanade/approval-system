package main

import (
	"main/config"
	"main/repository"
	"main/router"

	repositoryItem "main/repository/item"

	serviceItem "main/service/item"

	controllerItem "main/controller/item"
)

var (
	configManager config.ConfigManager = config.NewEnvConfigManager()
	database      repository.Database  = repository.NewDatabase(configManager)

	itemRepository repositoryItem.ItemRepository = repositoryItem.NewItemRepository(database)
	itemService    serviceItem.ItemService       = serviceItem.NewItemService(itemRepository, configManager)
	itemController controllerItem.ItemController = controllerItem.NewItemController(itemService)

	httpRouter router.Router = router.NewMuxRouter()
)
