package router

import "net/http"

type Router interface {
	GET(uri string, f func(resp http.ResponseWriter, req *http.Request))
	POST(uri string, f func(resp http.ResponseWriter, req *http.Request))
	SERVE(port string)
}
