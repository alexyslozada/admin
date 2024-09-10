package http

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"gitlab.com/EDteam/workshop-ai-2024/admin/domain"
	"gitlab.com/EDteam/workshop-ai-2024/admin/ports"
)

type LoginHandler struct {
	loginUseCase ports.LoginUseCase
}

func NewLoginHandler(loginUseCase ports.LoginUseCase) LoginHandler {
	return LoginHandler{loginUseCase: loginUseCase}
}

func (lh LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Bind body request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Unmarshal body request
	var loginRequest domain.LoginRequest
	err = json.Unmarshal(body, &loginRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call use case
	token, err := lh.loginUseCase.Login(loginRequest.Email, loginRequest.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Marshal response
	loginResponse := domain.LoginResponse{Token: token}
	response, err := json.Marshal(loginResponse)
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

func (lh LoginHandler) ValidateJWT(w http.ResponseWriter, r *http.Request) {
	// Bind body request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Unmarshal body request
	var validateJWTRequest domain.ValidateJWTRequest
	err = json.Unmarshal(body, &validateJWTRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call use case
	isValid, err := lh.loginUseCase.ValidateToken(validateJWTRequest.Token)
	msg := ""
	if err != nil {
		msg = err.Error()
	}

	// Marshal response
	validateJWTResponse := domain.ValidateJWTResponse{IsValid: isValid, Message: msg}
	response, err := json.Marshal(validateJWTResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	if !isValid {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	_, err = w.Write(response)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
