package di

import "net/http"

func Router() *http.ServeMux {

	loginHandler := InitLogin()
	clientHandler := InitClientHandler()
	saleHandler := InitSaleHandler()
	middleware := InitMiddleware()

	EDrouter := http.ServeMux{}
	EDrouter.HandleFunc("POST /v1/login", loginHandler.Login)
	EDrouter.HandleFunc("POST /v1/login/validate-jwt", loginHandler.ValidateJWT)
	EDrouter.HandleFunc("POST /v1/clients", middleware.Auth(clientHandler.Create))
	EDrouter.HandleFunc("GET /v1/clients", middleware.Auth(clientHandler.FindAll))
	EDrouter.HandleFunc("POST /v1/sales", middleware.Auth(saleHandler.Create))
	EDrouter.HandleFunc("GET /v1/sales", middleware.Auth(saleHandler.FindAll))

	return &EDrouter
}
