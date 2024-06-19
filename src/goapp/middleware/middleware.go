package middleware

import (
	"main/pkg/session"
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func AzureAuth() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			isAuth := session.IsAuthenticated(w, r)
			if isAuth {
				f(w, r)
			}
		}
	}
}

func ManagedIdentityAuth() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			session.IsGuidAuthenticated(w, r)
			f(w, r)
		}
	}
}

// Chain applies middlewares to a http.HandlerFunc
func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}
