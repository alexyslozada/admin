package http

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"gitlab.com/EDteam/workshop-ai-2024/admin/domain"
	"gitlab.com/EDteam/workshop-ai-2024/admin/ports/app"
)

type AIHandler struct {
	useCase app.AI
}

func NewAIHandler(useCase app.AI) AIHandler {
	return AIHandler{useCase: useCase}
}

// CreateThread creates a new thread using the AI use case and writes the response to the provided writer.
func (ah AIHandler) CreateThread(w http.ResponseWriter, r *http.Request) {
	// Call use case
	threadID, err := ah.useCase.CreateThread(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]string{"threadID": threadID.String()}
	responseRaw, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(responseRaw)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// CreateMessage creates a new message using the AI use case and writes the response to the provided writer.
func (ah AIHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	// Bind body request
	var request domain.AIMessageRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse uuid thread ID
	threadID, err := uuid.Parse(request.ThreadID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call use case
	response, err := ah.useCase.CreateMessage(r.Context(), threadID, request.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	responseRaw, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(responseRaw)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
