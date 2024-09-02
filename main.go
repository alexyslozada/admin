package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"gitlab.com/EDteam/workshop-ai-2024/admin/infrastructure/db"
	"gitlab.com/EDteam/workshop-ai-2024/admin/infrastructure/di"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file, err: %v", err)
	}

	// Migration
	db.RunMigration()

	loginHandler := di.InitLogin()

	EDrouter := http.ServeMux{}
	EDrouter.HandleFunc("POST /v1/login", loginHandler.Login)

	// Serve and listen
	log.Printf("Server running on port 8080")
	err = http.ListenAndServe(":8080", &EDrouter)
	if err != nil {
		panic(err)
	}
}
