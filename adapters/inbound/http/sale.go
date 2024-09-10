package http

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"gitlab.com/EDteam/workshop-ai-2024/admin/domain"
	"gitlab.com/EDteam/workshop-ai-2024/admin/ports"
)

type SaleHandler struct {
	useCase ports.GenericUseCase[domain.Sale]
}

func NewSaleHandler(useCase ports.GenericUseCase[domain.Sale]) SaleHandler {
	return SaleHandler{useCase: useCase}
}

func (h SaleHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Bind body request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Unmarshal body request
	var sale domain.Sale
	err = json.Unmarshal(body, &sale)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call use case
	err = h.useCase.Create(&sale)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Marshal response
	response, err := json.Marshal(sale)
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

func (s SaleHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	// Call use case
	sales, err := s.useCase.FindAll()
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
