package handlers

import (
	"context"
	"github.com/radiologist-ai/web-app/internal/domain"
	"net/http"
)

func (h *Handlers) WithHTMLResponse(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		handler(w, r)
	}
}

func (h *Handlers) WithCurrentUser(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			token string
			email string
			user  domain.UserRepoModel
			ok    bool
		)
		cookie, err := r.Cookie(domain.AuthTokenCookieKey)
		if err != nil {
			goto handle
		}

		token = cookie.Value
		if token == "" {
			goto handle
		}
		email, err = h.users.ValidateToken(h.secret, token)
		if err != nil {
			goto handle
		}
		if user, ok, err = h.users.GetByEmail(r.Context(), email); err != nil {
			h.logger.Error().Err(err).Str("email", email).Str("token", token).Msg("error")
			goto handle
		} else if !ok {
			h.logger.Warn().Bool("userExists", ok).Str("email", email).Str("token", token).Msg("user not found")
			goto handle
		}
		r = r.WithContext(context.WithValue(r.Context(), domain.CurrentUserCtxKey, user))

	handle:
		handler(w, r)
	}
}

func (h *Handlers) AnonymousRequired(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := GetCurrentUser(r.Context()); ok {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		handler(w, r)
	}
}

func (h *Handlers) AuthRequired(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := GetCurrentUser(r.Context()); !ok {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		handler(w, r)
	}
}
