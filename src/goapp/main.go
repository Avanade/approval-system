package main

import (
	"log"
	session "main/pkg/session"
	rtApprovals "main/routes/pages/approvals"
	"strconv"
	"time"

	ev "main/pkg/envvar"

	"github.com/joho/godotenv"
)

func main() {
	// Set environment variables
	err := godotenv.Load()
	if err != nil {
		log.Print(err.Error())
	}

	// Create session and GitHubClient
	session.InitializeSession()

	go checkFailedCallbacks()

	setPageRoutes()
	setApiRoutes()
	setUtilityRoutes()

	serve()
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
