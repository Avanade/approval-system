package main

import (
	router "main/http"
	m "main/middleware"
	ev "main/pkg/envvar"
	rtApi "main/routes/apis"
	rtAzure "main/routes/login/azure"
	rtPages "main/routes/pages"
	rtApprovals "main/routes/pages/approvals"
)

var (
	httpRouter router.Router = router.NewMuxRouter()
)

func setPageRoutes() {
	httpRouter.GET("/", m.Chain(rtApprovals.MyRequestsHandler, m.AzureAuth()))
	httpRouter.GET("/myapprovals", m.Chain(rtApprovals.MyApprovalsHandler, m.AzureAuth()))
	httpRouter.GET("/response/{appGuid}/{appModuleGuid}/{itemGuid}/{isApproved}", m.Chain(rtApprovals.ResponseHandler, m.AzureAuth()))
	httpRouter.GET("/responsereassigned/{appGuid}/{appModuleGuid}/{itemGuid}/{isApproved}/{ApproveText}/{RejectText}", m.Chain(rtApprovals.ResponseReassignedeHandler, m.AzureAuth()))

	httpRouter.GET("/loginredirect", rtPages.LoginRedirectHandler)
	httpRouter.GET("/login/azure", rtAzure.LoginHandler)
	httpRouter.GET("/login/azure/callback", rtAzure.CallbackHandler)
	httpRouter.GET("/logout/azure", rtAzure.LogoutHandler)
}

func setApiRoutes() {
	httpRouter.GET("/api/request/types", m.Chain(rtApi.GetRequestTypes, m.AzureAuth()))
	httpRouter.POST("/api/request", rtApprovals.ApprovalRequestHandler)
	httpRouter.POST("/api/process", rtApprovals.ProcessResponseHandler)
	httpRouter.GET("/api/items/type/{type:[0-2]+}/status/{status:[0-3]+}", m.Chain(rtApi.GetItems, m.AzureAuth()))
	httpRouter.GET("/api/search/users/{search}", m.Chain(rtApi.SearchUserFromActiveDirectory, m.AzureAuth()))
	httpRouter.GET("/api/responsereassignedapi/{itemGuid}/{approver}/{ApplicationId}/{ApplicationModuleId}/{itemId}/{ApproveText}/{RejectText}", m.Chain(rtApprovals.ReAssignApproverHandler, m.AzureAuth()))
}

func setUtilityRoutes() {
	httpRouter.GET("/utility/fillout-approvalrequest-approvers", m.Chain(rtApi.FillOutApprovalRequestApprovers, m.ManagedIdentityAuth()))
}

func serve() {
	port := ev.GetEnvVar("PORT", "8080")
	httpRouter.SERVE(port)
}
