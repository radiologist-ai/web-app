package handlers

import (
	"errors"
	"github.com/radiologist-ai/web-app/internal/domain"
	"github.com/radiologist-ai/web-app/internal/domain/customerrors"
	"github.com/radiologist-ai/web-app/internal/views"
	"net/http"
)

func (h *Handlers) PostPatientHandler(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := GetCurrentUser(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var form domain.PatientRepoModel
	form.Name = r.FormValue("name")
	if form.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return // TODO render smth
	}
	form.PatientIdentifier = r.FormValue("identifier")

	patient, err := h.patients.CreatePatient(r.Context(), *currentUser, form)
	if err != nil {
		h.logger.Error().Err(err).Msg("error creating patient")
		if errors.Is(err, customerrors.AccessError) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		http.Redirect(w, r, "/internal_server_error", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/patients/"+patient.ID.String()+"/reports", http.StatusFound)
	return
}

func (h *Handlers) GetHomeHandler(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := GetCurrentUser(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if currentUser.IsDoctor {
		patients, err := h.patients.GetAll(r.Context(), *currentUser)
		if err != nil {
			h.logger.Error().Err(err).Msg("error getting patients")
			http.Redirect(w, r, "/internal_server_error", http.StatusFound)
			return
		}
		err = views.Layout(views.Home(patients), "My Patients").Render(r.Context(), w)
		if err != nil {
			h.logger.Error().Err(err).Msg("error rendering layout")
			http.Redirect(w, r, "/internal_server_error", http.StatusFound)
		}
		return
	} else {
		err := views.Layout(views.Home(nil), "Home").Render(r.Context(), w)
		if err != nil {
			h.logger.Error().Err(err).Msg("error rendering layout")
			http.Redirect(w, r, "/internal_server_error", http.StatusFound)
		}
		return
	}
}
