package http

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"gitlab.com/EDteam/workshop-ai-2024/admin/domain"
	"gitlab.com/EDteam/workshop-ai-2024/admin/internal/urler"
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

func (h SaleHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	// Get filters
	filters := urler.ParseQueryParams(r.URL.Query())

	// Call use case
	sales, err := h.useCase.FindAll(filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	salesDTO := make([]domain.SaleResponse, 0, len(sales))
	for _, sale := range sales {
		salesDTO = append(salesDTO, domain.SaleResponse{
			ID:             sale.ID,
			Product:        sale.Product,
			ClientID:       sale.ClientID,
			Client:         sale.Client,
			DateInvoice:    time.Unix(sale.DateInvoice, 0),
			Amount:         sale.Amount,
			IsSubscription: sale.IsSubscription,
			Months:         sale.Months,
			CreatedAt:      time.Unix(sale.CreatedAt, 0),
			UpdatedAt:      time.Unix(sale.UpdatedAt, 0),
			DeletedAt:      time.Unix(sale.DeletedAt, 0),
		})
	}

	time.Now().Unix()
	// Marshal response
	response, err := json.Marshal(salesDTO)
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
