package handlers

import "net/http"

func (h *Handlers) WithHTMLResponse(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		handler(w, r)
	}
}
