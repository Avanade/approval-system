package middleware

import "net/http"

type Middleware interface {
	AzureAuth() MiddlewareFunc
	ManagedIdentityAuth() MiddlewareFunc
	Chain(f http.HandlerFunc, middlewares ...MiddlewareFunc) http.HandlerFunc
}
