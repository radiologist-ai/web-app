package http

import (
	"errors"
	"github.com/radiologist-ai/web-app/internal/app/http/handlers"
	"net/http"
)

func NewRouter(handlers *handlers.Handlers) (*http.ServeMux, error) {
	if handlers == nil {
		return nil, errors.New("handler is nil")
	}
	mux := http.NewServeMux()

	// auth
	mux.HandleFunc("GET /", handlers.WithHTMLResponse(handlers.HandleIndex))
	mux.HandleFunc("GET /register", handlers.WithHTMLResponse(handlers.HandleRegister))
	return mux, nil
}
