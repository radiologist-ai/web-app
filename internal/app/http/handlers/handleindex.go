package handlers

import (
	"github.com/radiologist-ai/web-app/internal/views"
	"net/http"
)

func (h *Handlers) HandleIndex(w http.ResponseWriter, r *http.Request) {
	if err := views.Layout(views.Index(), "Radiologist AI").Render(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
