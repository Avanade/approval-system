package ipdisclosure

import "net/http"

type IpDisclosureController interface {
	InsertIPDisclosureRequest(w http.ResponseWriter, r *http.Request)
	UpdateResponse(w http.ResponseWriter, r *http.Request)
}

type IpDisclosurePageController interface {
	IpDisclosureRequest(w http.ResponseWriter, r *http.Request)
}
