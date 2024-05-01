package http

import (
	"errors"
	"github.com/a-h/templ"
	"github.com/radiologist-ai/web-app/internal/app/http/handlers"
	"github.com/radiologist-ai/web-app/internal/views"
	"net/http"
)

func NewRouter(handlers *handlers.Handlers) (*http.ServeMux, error) {
	if handlers == nil {
		return nil, errors.New("handler is nil")
	}
	mux := http.NewServeMux()

	// auth
	mux.HandleFunc("GET /{$}", handlers.WithHTMLResponse(handlers.HandleIndex))
	mux.HandleFunc("GET /register", handlers.WithHTMLResponse(handlers.HandleRegister))
	mux.HandleFunc("POST /register", handlers.PostRegister)
	mux.HandleFunc("POST /validate/email", handlers.WithHTMLResponse(handlers.ValidateEmail))
	mux.HandleFunc("POST /validate/password", handlers.WithHTMLResponse(handlers.ValidatePassword))

	// technical
	mux.HandleFunc("GET /internal_server_error", handlers.WithHTMLResponse(templ.Handler(views.Layout(views.InternalError(), "Internal Error")).ServeHTTP))

	mux.HandleFunc("GET /", handlers.WithHTMLResponse(templ.Handler(views.Layout(views.NotFound(), "404")).ServeHTTP))
	return mux, nil
}
