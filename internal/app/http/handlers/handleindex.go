package handlers

import (
	"ai-radiologist/internal/views"
	"net/http"
)

func (h *Handlers) HandleIndex(w http.ResponseWriter, r *http.Request) {
	if err := views.Index().Render(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
