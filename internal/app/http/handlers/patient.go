package handlers

import (
	"errors"
	"github.com/radiologist-ai/web-app/internal/domain"
	"github.com/radiologist-ai/web-app/internal/domain/customerrors"
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
		return
	}
	form.PatientIdentifier = r.FormValue("identifier")

	patient, err := h.patients.CreatePatient(r.Context(), *currentUser, form)
	if err != nil {
		h.logger.Error().Err(err).Msg("error creating patient")
		if errors.Is(err, customerrors.AccessError) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/patients/"+patient.ID.String()+"/reports", http.StatusFound)
	return
}

// TODO get my patients list
