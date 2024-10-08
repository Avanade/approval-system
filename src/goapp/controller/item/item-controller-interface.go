package item

import "net/http"

type ItemController interface {
	GetItems(w http.ResponseWriter, r *http.Request)
	CreateItem(w http.ResponseWriter, r *http.Request)
	ProcessResponse(w http.ResponseWriter, r *http.Request)
	ReassignItem(w http.ResponseWriter, r *http.Request)
}

type ItemPageController interface {
	MyRequests(w http.ResponseWriter, r *http.Request)
	MyApprovals(w http.ResponseWriter, r *http.Request)
}
