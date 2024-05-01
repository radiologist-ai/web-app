package handlers

import (
	"errors"
	"github.com/radiologist-ai/web-app/internal/domain"
	"github.com/radiologist-ai/web-app/internal/domain/customerrors"
	"github.com/radiologist-ai/web-app/internal/views"
	"net/http"
	"net/mail"
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

func (h *Handlers) PostRegister(w http.ResponseWriter, r *http.Request) {
	var form domain.UserForm
	form.Email = r.FormValue("email")
	form.Password = r.FormValue("password")
	form.FirstName = r.FormValue("firstName")
	form.LastName = r.FormValue("lastName")
	form.IsDoctor = r.FormValue("isDoctor") == "on"

	if _, err := h.users.CreateOne(r.Context(), form); err != nil {
		if errors.Is(err, customerrors.ValidationError) {
			errTxt := ValidationErrorToResponseText(err)
			if err := views.Layout(views.RegistrationFormBad(errTxt), "Radiologist AI").Render(r.Context(), w); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		http.Redirect(w, r, "/internal_server_error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/internal_server_error", http.StatusCreated)
	return
}

func (h *Handlers) ValidateEmail(w http.ResponseWriter, r *http.Request) {
	var email = r.FormValue("email")
	if email == "" {
		if err := views.EmailInput("", email).Render(r.Context(), w); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	if _, err := mail.ParseAddress(email); err != nil {
		if err = views.EmailInput("is-invalid", email).Render(r.Context(), w); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	if err := views.EmailInput("is-valid", email).Render(r.Context(), w); err != nil {
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
