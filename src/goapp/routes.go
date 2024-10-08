package main

import (
	m "main/middleware"
	ev "main/pkg/envvar"
	rtAzure "main/routes/login/azure"
	rtPages "main/routes/pages"
	rtApprovals "main/routes/pages/approvals"
)

func setPageRoutes() {
	httpRouter.GET("/", m.Chain(ctrl.ItemPage.MyRequests, m.AzureAuth()))
	httpRouter.GET("/myapprovals", m.Chain(rtApprovals.MyApprovalsHandler, m.AzureAuth()))
	httpRouter.GET("/response/{appGuid}/{appModuleGuid}/{itemGuid}/{isApproved}", m.Chain(rtApprovals.ResponseHandler, m.AzureAuth()))
	httpRouter.GET("/responsereassigned/{appGuid}/{appModuleGuid}/{itemGuid}/{isApproved}/{ApproveText}/{RejectText}", m.Chain(rtApprovals.ResponseReassignedeHandler, m.AzureAuth()))

	httpRouter.GET("/loginredirect", rtPages.LoginRedirectHandler)
	httpRouter.GET("/login/azure", rtAzure.LoginHandler)
	httpRouter.GET("/login/azure/callback", rtAzure.CallbackHandler)
	httpRouter.GET("/logout/azure", rtAzure.LogoutHandler)
}

func setApiRoutes() {
	httpRouter.GET("/api/request/types", m.Chain(ctrl.ApplicationModule.GetRequestTypes, m.AzureAuth()))
	httpRouter.POST("/api/request", ctrl.Item.CreateItem)
	httpRouter.POST("/api/process", ctrl.Item.ProcessResponse)
	httpRouter.GET("/api/items/type/{type:[0-2]+}/status/{status:[0-3]+}", m.Chain(ctrl.Item.GetItems, m.AzureAuth()))
	httpRouter.GET("/api/search/users/{search}", m.Chain(ctrl.User.SearchUserFromActiveDirectory, m.AzureAuth()))
	httpRouter.GET("/api/responsereassignedapi/{itemGuid}/{approver}/{ApplicationId}/{ApplicationModuleId}/{ApproveText}/{RejectText}", m.Chain(rtApprovals.ReAssignApproverHandler, m.AzureAuth()))
}

func serve() {
	port := ev.GetEnvVar("PORT", "8080")
	httpRouter.SERVE(port)
}
