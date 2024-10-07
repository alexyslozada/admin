package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"gitlab.com/EDteam/workshop-ai-2024/admin/infrastructure/db"
	"gitlab.com/EDteam/workshop-ai-2024/admin/infrastructure/di"
)

func main() {
	// Config log to show the line of the log
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// Config log to standard output
	log.SetOutput(os.Stdout)

	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file, err: %v", err)
	}

	// Migration
	db.RunMigration()

	router := di.Router()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Serve and listen
	log.Printf("Server running on port 8080")
	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		panic(err)
	}
}
