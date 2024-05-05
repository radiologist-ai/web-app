package handlers

import (
	"context"
	"errors"
	"github.com/radiologist-ai/web-app/internal/domain"
	"github.com/radiologist-ai/web-app/internal/domain/customerrors"
	"github.com/radiologist-ai/web-app/internal/views"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/mail"
	"time"
)

func (h *Handlers) PostLogout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{Name: domain.AuthTokenCookieKey, Value: "", Expires: time.Time{}}
	http.SetCookie(w, &cookie)
	ctx := context.WithValue(r.Context(), domain.CurrentUserCtxKey, nil)
	r = r.WithContext(ctx)
	if err := views.Nav(nil).Render(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	return
}

func (h *Handlers) PostLogin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, ok, err := h.users.GetByEmail(r.Context(), email)
	if err != nil {
		http.Redirect(w, r, "/internal_server_error", http.StatusFound)
		return
	}
	if !ok {
		if err := views.Layout(views.LoginFormUserDoesntExist(), "Radiologist AI").Render(r.Context(), w); err != nil {
			http.Redirect(w, r, "/internal_server_error", http.StatusFound)
		}
		return
	}
	if bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)) != nil {
		if err := views.Layout(views.LoginFormWrongPassword(), "Radiologist AI").Render(r.Context(), w); err != nil {
			http.Redirect(w, r, "/internal_server_error", http.StatusFound)
		}
	}

	token, err := h.users.GenerateToken(h.secret, user.Email)
	if err != nil {
		h.logger.Error().Err(err).Any("user", user).Msg("Error generating token")
		http.Redirect(w, r, "/internal_server_error", http.StatusFound)
		return
	}

	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: domain.AuthTokenCookieKey, Value: token, Expires: expiration}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/home", http.StatusFound)
	return
}

func (h *Handlers) PostRegister(w http.ResponseWriter, r *http.Request) {
	var form domain.UserForm
	form.Email = r.FormValue("email")
	form.Password = r.FormValue("password")
	form.FirstName = r.FormValue("firstName")
	form.LastName = r.FormValue("lastName")
	form.IsDoctor = r.FormValue("isDoctor") == "on"

	user, err := h.users.CreateOne(r.Context(), form)
	if err != nil {
		if errors.Is(err, customerrors.ValidationError) {
			errTxt := ValidationErrorToResponseText(err)
			if err := views.Layout(views.RegistrationFormBad(errTxt), "Radiologist AI").Render(r.Context(), w); err != nil {
				http.Redirect(w, r, "/internal_server_error", http.StatusFound)
			}
			return
		}
		http.Redirect(w, r, "/internal_server_error", http.StatusFound)
		return
	}
	token, err := h.users.GenerateToken(h.secret, user.Email)
	if err != nil {
		h.logger.Error().Err(err).Any("user", user).Msg("Error generating token")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: domain.AuthTokenCookieKey, Value: token, Expires: expiration}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/home", http.StatusFound)
	return
}

func (h *Handlers) ValidateEmail(w http.ResponseWriter, r *http.Request) {
	var email = r.FormValue("email")
	if email == "" {
		if err := views.EmailInput("", email, "Invalid Email").Render(r.Context(), w); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	if _, err := mail.ParseAddress(email); err != nil {
		if err = views.EmailInput("is-invalid", email, "Invalid Email").Render(r.Context(), w); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	if _, ok, _ := h.users.GetByEmail(r.Context(), email); ok {
		if err := views.EmailInput("is-invalid", email, "User with same email already exists").Render(r.Context(), w); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	if err := views.EmailInput("is-valid", email, "Invalid Email").Render(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handlers) ValidatePassword(w http.ResponseWriter, r *http.Request) {
	var password = r.FormValue("password")
	if password == "" {
		if err := views.PasswordInput("", password, "").Render(r.Context(), w); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	if err := h.users.ValidatePassword(password); err != nil {
		feedback := "Invalid password" // TODO may be return status 500?
		if errors.Is(err, customerrors.ValidationErrorPassword) {
			feedback = ValidationErrorToResponseText(err)
		}
		if err := views.PasswordInput("is-invalid", password, feedback).Render(r.Context(), w); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	if err := views.PasswordInput("is-valid", password, "").Render(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
