package appmodule

import "net/http"

type ApplicationModuleController interface {
	GetRequestTypes(w http.ResponseWriter, r *http.Request)
}
