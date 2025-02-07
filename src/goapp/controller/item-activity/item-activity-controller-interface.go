package itemActivity

import (
	"net/http"
)

type ItemActivityController interface {
	InsertItemActivity(w http.ResponseWriter, r *http.Request)
	GetItemActivity(w http.ResponseWriter, r *http.Request)
}
