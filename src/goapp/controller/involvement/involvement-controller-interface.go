package involvement

import (
	"net/http"
)

type InvolvementController interface {
	GetInvolvementList(w http.ResponseWriter, r *http.Request)
}
