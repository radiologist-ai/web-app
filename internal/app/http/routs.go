package http

import (
	"ai-radiologist/internal/app/http/handlers"
	"errors"
	"net/http"
)

func NewRouter(handlers *handlers.Handlers) (*http.ServeMux, error) {
	if handlers == nil {
		return nil, errors.New("handler is nil")
	}
	mux := http.NewServeMux()

	// auth
	mux.HandleFunc("GET /", handlers.WithHTMLResponse(handlers.HandleIndex))

	return mux, nil
}
