package authentication

import (
	"net/http"
)

type AuthenticationPageController interface {
	AuthenticationFailedHandler(w http.ResponseWriter, r *http.Request)
	AuthenticationInProgressHandler(w http.ResponseWriter, r *http.Request)
	AuthenticationSuccessfulHandler(w http.ResponseWriter, r *http.Request)
	CallbackHandler(w http.ResponseWriter, r *http.Request)
	LoginHandler(w http.ResponseWriter, r *http.Request)
	LoginRedirectHandler(w http.ResponseWriter, r *http.Request)
	LogoutHandler(w http.ResponseWriter, r *http.Request)
}
