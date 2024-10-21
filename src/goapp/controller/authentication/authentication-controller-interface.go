package authentication

import (
	"net/http"
)

type AuthenticationPageController interface {
	LoginRedirectHandler(w http.ResponseWriter, r *http.Request)
}
