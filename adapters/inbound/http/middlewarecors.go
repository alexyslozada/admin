package http

import "net/http"

type EDmux struct {
	mux            *http.ServeMux
	allowedOrigins map[string]struct{}
}

func NewEDmux(allowedOrigins map[string]struct{}) EDmux {
	return EDmux{mux: http.NewServeMux(), allowedOrigins: allowedOrigins}
}

func (e *EDmux) HandleFunc(pattern string, handler http.HandlerFunc) {
	e.mux.HandleFunc(pattern, handler)
}

func (e *EDmux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")

	_, ok := e.allowedOrigins[origin]
	if ok {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	}

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	e.mux.ServeHTTP(w, r)
}
