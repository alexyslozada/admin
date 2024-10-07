package http

import (
	"encoding/json"
	"log"
	"net/http"

	"gitlab.com/EDteam/workshop-ai-2024/admin/domain"
	"gitlab.com/EDteam/workshop-ai-2024/admin/internal/urler"
	"gitlab.com/EDteam/workshop-ai-2024/admin/ports/app"
)

type SaleSummarizedHandler struct {
	useCase app.GenericUseCase[domain.SaleSummarized]
}

func NewSaleSummarizedHandler(useCase app.GenericUseCase[domain.SaleSummarized]) SaleSummarizedHandler {
	return SaleSummarizedHandler{useCase: useCase}
}

func (h SaleSummarizedHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	// Get filters
	filters := urler.ParseQueryParams(r.URL.Query())

	// Call use case
	sales, err := h.useCase.FindAll(filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Marshal response
	response, err := json.Marshal(sales)
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
