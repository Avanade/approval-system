package main

import (
	"fmt"
	"log"
	session "main/pkg/session"
	rtApi "main/routes/apis"
	rtAzure "main/routes/login/azure"
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
		IsDevelopment:         false,
	})

	// Set environment variables
	err := godotenv.Load()
	if err != nil {
		log.Print(err.Error())
	}

	// Create session and GitHubClient
	session.InitializeSession()

	mux := mux.NewRouter()
	mux.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	mux.Handle("/", loadAzAuthPage(rtApprovals.MyRequestsHandler))
	mux.Handle("/myapprovals", loadAzAuthPage(rtApprovals.MyApprovalsHandler))
	mux.Handle("/response/{appGuid}/{appModuleGuid}/{itemGuid}/{isApproved}", loadAzAuthPage(rtApprovals.ResponseHandler))
	mux.Handle("/responseReassigned/{appGuid}/{appModuleGuid}/{itemGuid}/{isApproved}", loadAzAuthPage(rtApprovals.ResponseReassignedeHandler))
	mux.HandleFunc("/loginredirect", rtPages.LoginRedirectHandler).Methods("GET")
	mux.HandleFunc("/login/azure", rtAzure.LoginHandler)
	mux.HandleFunc("/login/azure/callback", rtAzure.CallbackHandler)
	mux.HandleFunc("/logout/azure", rtAzure.LogoutHandler)
	mux.HandleFunc("/request", rtApprovals.ApprovalRequestHandler)
	mux.HandleFunc("/process", rtApprovals.ProcessResponseHandler)
	muxApi := mux.PathPrefix("/api").Subrouter()
	muxApi.Handle("/items/type/{type:[0-2]+}/status/{status:[0-3]+}", loadAzAuthPage(rtApi.GetItems))
	muxApi.Handle("/search/users/{search}", loadAzAuthPage(rtApi.SearchUserFromActiveDirectory))
	muxApi.Handle("/responseReassignedAPI/{itemGuid}/{approver}", loadAzAuthPage(rtApprovals.ReAssignApproverHandler))

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
