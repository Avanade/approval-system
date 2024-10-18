package main

import (
	"log"
	session "main/pkg/session"

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

	go timedJobs.ReprocessFailedCallbacks()

	setPageRoutes()
	setApiRoutes()

	serve()
}
