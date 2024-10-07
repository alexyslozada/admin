package di

import (
	"log"
	"os"
	"strings"

	EDhttp "gitlab.com/EDteam/workshop-ai-2024/admin/adapters/inbound/http"
)

func Router() *EDhttp.EDmux {

	loginHandler := InitLogin()
	clientHandler := InitClientHandler()
	saleHandler := InitSaleHandler()
	saleSummarizedHandler := InitSaleSummarizedHandler()
	aiHandler := InitAIHandler()
	middleware := InitMiddleware()

	// CORS
	allowedDomains := os.Getenv("ALLOWED_DOMAINS")
	allowedDomainsList := strings.Split(allowedDomains, ",")
	allowedDomainsUnique := make(map[string]struct{}, len(allowedDomainsList))
	for _, domain := range allowedDomainsList {
		allowedDomainsUnique[domain] = struct{}{}
	}

	log.Printf("Allowed domains: %v", allowedDomainsList)

	EDrouter := EDhttp.NewEDmux(allowedDomainsUnique)
	EDrouter.HandleFunc("POST /v1/login", loginHandler.Login)
	EDrouter.HandleFunc("POST /v1/login/validate-jwt", loginHandler.ValidateJWT)
	EDrouter.HandleFunc("POST /v1/clients", middleware.Auth(clientHandler.Create))
	EDrouter.HandleFunc("GET /v1/clients", middleware.Auth(clientHandler.FindAll))
	EDrouter.HandleFunc("POST /v1/sales", middleware.Auth(saleHandler.Create))
	EDrouter.HandleFunc("GET /v1/sales", middleware.Auth(saleHandler.FindAll))
	EDrouter.HandleFunc("GET /v1/sales-summarized", middleware.Auth(saleSummarizedHandler.FindAll))
	EDrouter.HandleFunc("POST /v1/ai/thread", middleware.Auth(aiHandler.CreateThread))
	EDrouter.HandleFunc("POST /v1/ai/message", middleware.Auth(aiHandler.CreateMessage))

	return &EDrouter
}
