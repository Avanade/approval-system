package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Set environment variables
	err := godotenv.Load()
	if err != nil {
		log.Print(err.Error())
	}

	go timedJobs.ReprocessFailedCallbacks()

	setPageRoutes()
	setApiRoutes()

	serve()
}
