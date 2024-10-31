package router

import (
	"fmt"
	"net/http"
	"os"

	"main/config"
	"main/controller"
	"main/middleware"

	"github.com/gorilla/mux"
	"github.com/unrolled/secure"
)

type muxRouter struct {
	*controller.Controller
	Port string
	m    middleware.Middleware
}

func NewMuxRouter(c *controller.Controller, conf config.ConfigManager, m *middleware.Middleware) Router {
	return &muxRouter{
		Controller: c,
		Port:       conf.GetPort(),
		m:          *m,
	}
}

var (
	muxDispatcher = mux.NewRouter()
)

func (*muxRouter) GET(uri string, f func(resp http.ResponseWriter, req *http.Request)) {
	muxDispatcher.HandleFunc(uri, f).Methods("GET")
}

func (*muxRouter) POST(uri string, f func(resp http.ResponseWriter, req *http.Request)) {
	muxDispatcher.HandleFunc(uri, f).Methods("POST")
}

func (*muxRouter) PUT(uri string, f func(resp http.ResponseWriter, req *http.Request)) {
	muxDispatcher.HandleFunc(uri, f).Methods("PUT")
}

func (*muxRouter) DELETE(uri string, f func(resp http.ResponseWriter, req *http.Request)) {
	muxDispatcher.HandleFunc(uri, f).Methods("DELETE")
}

func (r *muxRouter) SERVE() {
	secureOptions := secure.Options{
		SSLRedirect:           true,                                            // Strict-Transport-Security
		SSLHost:               os.Getenv("SSL_HOST"),                           // Strict-Transport-Security
		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"}, // Strict-Transport-Security
		FrameDeny:             true,
		ContentTypeNosniff:    true, // X-Content-Type-Options
		BrowserXssFilter:      true,
		ReferrerPolicy:        "strict-origin", // Referrer-Policy
		ContentSecurityPolicy: os.Getenv("CONTENT_SECURITY_POLICY"),
		PermissionsPolicy:     "fullscreen=(), geolocation=()", // Permissions-Policy
		STSSeconds:            31536000,                        // Strict-Transport-Security
		STSIncludeSubdomains:  true,                            // Strict-Transport-Security
		IsDevelopment:         os.Getenv("IS_DEVELOPMENT") == "true",
	}

	secureMiddleware := secure.New(secureOptions)
	muxDispatcher.Use(secureMiddleware.Handler)
	http.Handle("/", muxDispatcher)

	muxDispatcher.NotFoundHandler = http.HandlerFunc(r.m.Chain(r.Controller.Fallback.NotFound, r.m.AzureAuth()))
	muxDispatcher.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))

	fmt.Printf("Mux HTTP server running on port %v", r.Port)
	http.ListenAndServe(fmt.Sprintf(":%v", r.Port), muxDispatcher)
}
