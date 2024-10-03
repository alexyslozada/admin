package http

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"gitlab.com/EDteam/workshop-ai-2024/admin/domain"
	"gitlab.com/EDteam/workshop-ai-2024/admin/internal/urler"
	"gitlab.com/EDteam/workshop-ai-2024/admin/ports/app"
)

type ClientHandler struct {
	useCase app.GenericUseCase[domain.Client]
}

func NewClientHandler(useCase app.GenericUseCase[domain.Client]) ClientHandler {
	return ClientHandler{useCase: useCase}
}

func (ch ClientHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Bind body request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Unmarshal body request
	var client domain.Client
	err = json.Unmarshal(body, &client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call use case
	err = ch.useCase.Create(&client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Marshal response
	response, err := json.Marshal(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(response)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func (ch ClientHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	// Get filters
	filters := urler.ParseQueryParams(r.URL.Query())

	// Call use case
	clients, err := ch.useCase.FindAll(filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Marshal response
	response, err := json.Marshal(clients)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
