package authentication

import (
	"net/http"
)

type AuthenticationPageController interface {
	LoginHandler(w http.ResponseWriter, r *http.Request)
	LoginRedirectHandler(w http.ResponseWriter, r *http.Request)
}
