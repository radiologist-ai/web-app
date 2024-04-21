package handlers

import (
	"github.com/radiologist-ai/web-app/internal/views"
	"net/http"
)

func (h *Handlers) HandleRegister(w http.ResponseWriter, r *http.Request) {
	_, ok := GetCurrentUser(r.Context())
	if !ok {
		if err := views.Layout(views.RegistrationForm(), "Radiologist AI").Render(r.Context(), w); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	http.Redirect(w, r, "/", http.StatusOK)
}
