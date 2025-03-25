package middleware

import (
	"fmt"
	"main/service"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type middleware struct {
	*service.Service
}

func NewMiddleware(s *service.Service) Middleware {
	return &middleware{
		Service: s,
	}
}

func (m *middleware) AzureAuth() MiddlewareFunc {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Check if the request is for a response endpoint
			var url string
			url = fmt.Sprintf("%v", r.URL)
			if strings.HasPrefix(url, "/response/") {
				authReq := true
				params := mux.Vars(r)
				appModuleGuid := params["appModuleGuid"]

				authReq, err := m.Service.ApplicationModule.IsAuthRequired(appModuleGuid)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				if !authReq {
					f(w, r)
					return
				}
			}

			// Auth check for all other requests
			url = fmt.Sprintf("/loginredirect?redirect=%v", r.URL)
			isAuth, tokenExpired, err := m.Service.Authenticator.IsAuthenticated(r)
			if err != nil {
				c := http.Cookie{
					Name:   "auth-session",
					MaxAge: -1}
				http.SetCookie(w, &c)
				http.Redirect(w, r, url, http.StatusTemporaryRedirect)
				return
			}

			if isAuth {
				if tokenExpired {
					err := m.Service.Authenticator.RefreshToken(&w, r)
					if err != nil {
						c := http.Cookie{
							Name:   "auth-session",
							MaxAge: -1}
						http.SetCookie(w, &c)
						http.Redirect(w, r, "/logout", http.StatusTemporaryRedirect)
						return
					}
				}
				f(w, r)
			} else {
				http.Redirect(w, r, url, http.StatusTemporaryRedirect)
			}
		}
	}
}

func (m *middleware) ManagedIdentityAuth() MiddlewareFunc {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if m.Service.Authenticator.AccessTokenIsValid(r) {
				f(w, r)
			} else {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}
		}
	}
}

// Chain applies middlewares to a http.HandlerFunc
func (m *middleware) Chain(f http.HandlerFunc, middlewares ...MiddlewareFunc) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		f = middlewares[i](f)
	}
	return f
}
