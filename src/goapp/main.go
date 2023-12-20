package main

import (
	"fmt"
	"log"
	session "main/pkg/session"
	rtPages "main/routes/pages"
	rtApprovals "main/routes/pages/approvals"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/unrolled/secure"

	ev "main/pkg/envvar"

	"github.com/codegangsta/negroni"
	"github.com/joho/godotenv"
)

func main() {
	// Set environment variables
	err := godotenv.Load()
	if err != nil {
		log.Print(err.Error())
	}

	secureMiddleware := secure.New(secure.Options{
		SSLRedirect:           true,                                            // Strict-Transport-Security
		SSLHost:               os.Getenv("SSL_HOST"),                           // Strict-Transport-Security
		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"}, // Strict-Transport-Security
		FrameDeny:             true,                                            // X-FRAME-OPTIONS
		ContentTypeNosniff:    true,                                            // X-Content-Type-Options
		BrowserXssFilter:      true,
		ReferrerPolicy:        "strict-origin", // Referrer-Policy
		ContentSecurityPolicy: os.Getenv("CONTENT_SECURITY_POLICY"),
		PermissionsPolicy:     "fullscreen=(), geolocation=()", // Permissions-Policy
		STSSeconds:            31536000,                        // Strict-Transport-Security
		STSIncludeSubdomains:  true,                            // Strict-Transport-Security,
		IsDevelopment:         os.Getenv("IS_DEVELOPMENT") == "true",
	})

	// Create session and GitHubClient
	session.InitializeSession()

	mux := mux.NewRouter()

	setPageRoutes(mux)
	setApiRoutes(mux)
	setUtilityRoutes(mux)

	mux.NotFoundHandler = loadAzAuthPage(rtPages.NotFoundHandler)

	go checkFailedCallbacks()

	mux.Use(secureMiddleware.Handler)
	http.Handle("/", mux)

	port := ev.GetEnvVar("PORT", "8080")
	fmt.Printf("Now listening on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), mux))

}

// Verifies authentication before loading the page.
func loadAzAuthPage(f func(w http.ResponseWriter, r *http.Request)) *negroni.Negroni {
	return negroni.New(
		negroni.HandlerFunc(session.IsAuthenticated),
		negroni.Wrap(http.HandlerFunc(f)),
	)
}

func loadGuidAuthApi(f func(w http.ResponseWriter, r *http.Request)) *negroni.Negroni {
	return negroni.New(
		negroni.HandlerFunc(session.IsGuidAuthenticated),
		negroni.Wrap(http.HandlerFunc(f)),
	)
}

func checkFailedCallbacks() {
	// TIMER SERVICE
	freq := ev.GetEnvVar("CALLBACK_RETRY_FREQ", "15")
	freqInt, _ := strconv.ParseInt(freq, 0, 64)
	if freq > "0" {
		for range time.NewTicker(time.Duration(freqInt) * time.Minute).C {
			rtApprovals.ProcessFailedCallbacks()
		}
	}
}
