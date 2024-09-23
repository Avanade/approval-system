package item

import "net/http"

type ItemController interface {
	GetItems(w http.ResponseWriter, r *http.Request)
}
