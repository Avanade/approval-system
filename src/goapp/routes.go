package main

func setPageRoutes() {
	httpRouter.GET("/", m.Chain(ctrl.ItemPage.MyRequests, m.AzureAuth()))
	httpRouter.GET("/multiple-approvals", m.Chain(ctrl.ItemPage.MultipleApprovals, m.AzureAuth()))
	httpRouter.GET("/myapprovals", m.Chain(ctrl.ItemPage.MyApprovals, m.AzureAuth()))
	httpRouter.GET("/review", m.Chain(ctrl.ItemPage.ForReview, m.AzureAuth()))
	httpRouter.GET("/audit", m.Chain(ctrl.ItemPage.ForAudit, m.AzureAuth()))
	httpRouter.GET("/ipdisclosurerequest", m.Chain(ctrl.IPDisclourePage.IpDisclosureRequest, m.AzureAuth()))
	httpRouter.GET("/{action}/{appGuid}/{appModuleGuid}/{itemGuid}/{isApproved}", m.Chain(ctrl.ItemPage.RespondToItem, m.AzureAuth()))
	httpRouter.GET("/responsereassigned/{appGuid}/{appModuleGuid}/{itemGuid}/{isApproved}/{ApproveText}/{RejectText}", m.Chain(ctrl.ItemPage.ReassignApproval, m.AzureAuth()))

	httpRouter.GET("/loginredirect", ctrl.AuthenticationPage.LoginRedirectHandler)
	httpRouter.GET("/login/azure", ctrl.AuthenticationPage.LoginHandler)
	httpRouter.GET("/login/azure/callback", ctrl.AuthenticationPage.CallbackHandler)
	httpRouter.GET("/logout/azure", ctrl.AuthenticationPage.LogoutHandler)
}

func setApiRoutes() {
	httpRouter.GET("/api/request/types", m.Chain(ctrl.ApplicationModule.GetRequestTypes, m.AzureAuth()))
	httpRouter.POST("/api/request", m.Chain(ctrl.Item.CreateItem, m.ManagedIdentityAuth()))
	httpRouter.POST("/api/process", m.Chain(ctrl.Item.ProcessResponse, m.AzureAuth()))
	httpRouter.POST("/api/multiple/response/{response}", m.Chain(ctrl.Item.ProcessMultipleResponse, m.AzureAuth()))
	httpRouter.GET("/api/involvement", m.Chain(ctrl.Involvement.GetInvolvementList, m.AzureAuth()))
	httpRouter.GET("/api/approver/me/items", m.Chain(ctrl.Item.GetItemsByApprover, m.AzureAuth()))
	httpRouter.GET("/api/items/type/{type:[0-2]+}/status/{status:[0-3]+}", m.Chain(ctrl.Item.GetItems, m.AzureAuth()))
	httpRouter.GET("/api/items/forreview/{status}", m.Chain(ctrl.Item.GetItemsForReviewByConsultant, m.AzureAuth()))
	httpRouter.GET("/api/items/foraudit/{status}", m.Chain(ctrl.Item.GetItemsForAudit, m.AzureAuth()))
	httpRouter.GET("/api/search/users/{search}", m.Chain(ctrl.User.SearchUserFromActiveDirectory, m.AzureAuth()))
	httpRouter.GET("/api/responsereassignedapi/{itemGuid}/{approver}/{ApplicationId}/{ApplicationModuleId}/{ApproveText}/{RejectText}", m.Chain(ctrl.Item.ReassignItem, m.AzureAuth()))
	httpRouter.POST("/api/ipdisclosurerequest", m.Chain(ctrl.IPDisclosure.InsertIPDisclosureRequest, m.AzureAuth()))
	httpRouter.POST("/api/ipdisclosurerequest/consultlegal", m.Chain(ctrl.Item.ConsultLegal, m.AzureAuth()))
	httpRouter.POST("/api/activity", m.Chain(ctrl.ItemActivity.InsertItemActivity, m.AzureAuth()))
	httpRouter.GET("/api/activity/{id}", m.Chain(ctrl.ItemActivity.GetItemActivity, m.AzureAuth()))
}

func serve() {
	httpRouter.SERVE()
}
