package middleware

import "net/http"

type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc
