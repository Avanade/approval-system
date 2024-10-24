package main

func setPageRoutes() {
	httpRouter.GET("/", m.Chain(ctrl.ItemPage.MyRequests, m.AzureAuth()))
	httpRouter.GET("/myapprovals", m.Chain(ctrl.ItemPage.MyApprovals, m.AzureAuth()))
	httpRouter.GET("/response/{appGuid}/{appModuleGuid}/{itemGuid}/{isApproved}", m.Chain(ctrl.ItemPage.RespondToItem, m.AzureAuth()))
	httpRouter.GET("/responsereassigned/{appGuid}/{appModuleGuid}/{itemGuid}/{isApproved}/{ApproveText}/{RejectText}", m.Chain(ctrl.ItemPage.ReassignApproval, m.AzureAuth()))

	httpRouter.GET("/loginredirect", ctrl.AuthenticationPage.LoginRedirectHandler)
	httpRouter.GET("/login/azure", ctrl.AuthenticationPage.LoginHandler)
	httpRouter.GET("/login/azure/callback", ctrl.AuthenticationPage.CallbackHandler)
	httpRouter.GET("/logout/azure", ctrl.AuthenticationPage.LogoutHandler)
}

func setApiRoutes() {
	httpRouter.GET("/api/request/types", m.Chain(ctrl.ApplicationModule.GetRequestTypes, m.AzureAuth()))
	httpRouter.POST("/api/request", ctrl.Item.CreateItem)
	httpRouter.POST("/api/process", ctrl.Item.ProcessResponse)
	httpRouter.GET("/api/items/type/{type:[0-2]+}/status/{status:[0-3]+}", m.Chain(ctrl.Item.GetItems, m.AzureAuth()))
	httpRouter.GET("/api/search/users/{search}", m.Chain(ctrl.User.SearchUserFromActiveDirectory, m.AzureAuth()))
	httpRouter.GET("/api/responsereassignedapi/{itemGuid}/{approver}/{ApplicationId}/{ApplicationModuleId}/{ApproveText}/{RejectText}", m.Chain(ctrl.Item.ReassignItem, m.AzureAuth()))
}

func serve() {
	httpRouter.SERVE()
}
