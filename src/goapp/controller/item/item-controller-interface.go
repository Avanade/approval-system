package item

import "net/http"

type ItemController interface {
	ConsultLegal(w http.ResponseWriter, r *http.Request)
	CreateItem(w http.ResponseWriter, r *http.Request)
	GetItems(w http.ResponseWriter, r *http.Request)
	GetItemsByApprover(w http.ResponseWriter, r *http.Request)
	GetItemsForReviewByConsultant(w http.ResponseWriter, r *http.Request)
	ProcessResponse(w http.ResponseWriter, r *http.Request)
	ProcessMultipleResponse(w http.ResponseWriter, r *http.Request)
	ReassignItem(w http.ResponseWriter, r *http.Request)
}

type ItemPageController interface {
	ForReview(w http.ResponseWriter, r *http.Request)
	MyRequests(w http.ResponseWriter, r *http.Request)
	MyApprovals(w http.ResponseWriter, r *http.Request)
	MultipleApprovals(w http.ResponseWriter, r *http.Request)
	RespondToItem(w http.ResponseWriter, r *http.Request)
	ReassignApproval(w http.ResponseWriter, r *http.Request)
}
