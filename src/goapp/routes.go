package main

import (
	rtApi "main/routes/apis"
	rtAzure "main/routes/login/azure"
	rtPages "main/routes/pages"
	rtApprovals "main/routes/pages/approvals"
	"net/http"

	"github.com/gorilla/mux"
)

func setPageRoutes(mux *mux.Router) {
	mux.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	mux.Handle("/", loadAzAuthPage(rtApprovals.MyRequestsHandler))
	mux.Handle("/myapprovals", loadAzAuthPage(rtApprovals.MyApprovalsHandler))
	mux.Handle("/response/{appGuid}/{appModuleGuid}/{itemGuid}/{isApproved}", loadAzAuthPage(rtApprovals.ResponseHandler))
	mux.Handle("/responsereassigned/{appGuid}/{appModuleGuid}/{itemGuid}/{isApproved}/{ApproveText}/{RejectText}", loadAzAuthPage(rtApprovals.ResponseReassignedeHandler))

	mux.HandleFunc("/loginredirect", rtPages.LoginRedirectHandler)
	mux.HandleFunc("/login/azure", rtAzure.LoginHandler)
	mux.HandleFunc("/login/azure/callback", rtAzure.CallbackHandler)
	mux.HandleFunc("/logout/azure", rtAzure.LogoutHandler)
}

func setApiRoutes(mux *mux.Router) {
	muxApi := mux.PathPrefix("/api").Subrouter()
	muxApi.HandleFunc("/request", rtApprovals.ApprovalRequestHandler)
	muxApi.HandleFunc("/process", rtApprovals.ProcessResponseHandler)
	muxApi.Handle("/items/type/{type:[0-2]+}/status/{status:[0-3]+}", loadAzAuthPage(rtApi.GetItems))
	muxApi.Handle("/search/users/{search}", loadAzAuthPage(rtApi.SearchUserFromActiveDirectory))
	muxApi.Handle("/responsereassignedapi/{itemGuid}/{approver}/{ApplicationId}/{ApplicationModuleId}/{itemId}/{ApproveText}/{RejectText}", loadAzAuthPage(rtApprovals.ReAssignApproverHandler))
}

func setUtilityRoutes(mux *mux.Router) {
	muxUtility := mux.PathPrefix("/utility").Subrouter()
	muxUtility.Handle("/fillout-approvalrequest-approvers", loadGuidAuthApi(rtApi.FillOutApprovalRequestApprovers)).Methods("GET")
}
