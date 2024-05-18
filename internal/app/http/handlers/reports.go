package handlers

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/radiologist-ai/web-app/internal/domain/customerrors"
	"github.com/radiologist-ai/web-app/internal/views"
	"net/http"
)

func (h *Handlers) NewReportHandler(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := GetCurrentUser(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	patientID := r.PathValue("patientID")
	if patientID == "" {
		http.Redirect(w, r, "/not_found", http.StatusTemporaryRedirect)
		return
	}
	patientUUID, err := uuid.Parse(patientID)
	if err != nil {
		http.Redirect(w, r, "/not_found", http.StatusTemporaryRedirect)
		return
	}

	patient, err := h.patients.GetOne(r.Context(), *currentUser, patientUUID)
	if err != nil {
		if errors.Is(err, customerrors.NotFoundError) {
			w.WriteHeader(http.StatusNotFound)
		}
		http.Redirect(w, r, "/internal_server_error", http.StatusFound)
		return
	}
	if err := views.Layout(views.NewReportPage(patient), "New Report").Render(r.Context(), w); err != nil {
		http.Redirect(w, r, "/internal_server_error", http.StatusFound)
		return
	}

}

func (h *Handlers) PostNewReportHandler(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := GetCurrentUser(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	patientID := r.PathValue("patientID")
	if patientID == "" {
		http.Redirect(w, r, "/not_found", http.StatusTemporaryRedirect)
		return
	}
	patientUUID, err := uuid.Parse(patientID)
	if err != nil {
		http.Redirect(w, r, "/not_found", http.StatusTemporaryRedirect)
		return
	}

	patient, err := h.patients.GetOne(r.Context(), *currentUser, patientUUID)
	if err != nil {
		if errors.Is(err, customerrors.NotFoundError) {
			w.WriteHeader(http.StatusNotFound)
		}
		http.Redirect(w, r, "/internal_server_error", http.StatusFound)
		return
	}

	file, fileHeader, err := r.FormFile("xray")
	if err != nil || fileHeader == nil {
		h.logger.Error().Err(err).Any("fileHeader", fileHeader).Msg("r.FormFile(\"xray\")")
		http.Redirect(w, r, "/internal_server_error", http.StatusFound)
		return
	}

	report, err := h.rgen.GenerateReport(r.Context(), patient.ID, file, fileHeader.Size)
	if err != nil {
		h.logger.Error().Err(err).Msg("rgen.GenerateReport()")
		http.Redirect(w, r, "/internal_server_error", http.StatusFound)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/patients/%s/reports/%d", patient.ID.String(), report.ID), http.StatusFound)
	return
}
