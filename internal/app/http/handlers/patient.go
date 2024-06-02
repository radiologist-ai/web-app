package handlers

import (
	"errors"
	"github.com/google/uuid"
	"github.com/radiologist-ai/web-app/internal/domain"
	"github.com/radiologist-ai/web-app/internal/domain/customerrors"
	"github.com/radiologist-ai/web-app/internal/views"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

func (h *Handlers) GetMyAccountsHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info().Any("path", r.URL).Msg("request received GetMyAccountsHandler")
	currentUser, ok := GetCurrentUser(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	accounts, err := h.patients.GetAllByUser(r.Context(), *currentUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := views.ListOfAccounts(accounts).Render(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}

func (h *Handlers) PostLinkAccountHandler(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := GetCurrentUser(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var (
		comment string
		success bool
	)
	patientID := r.FormValue("patientID")
	if patientID == "" {
		comment = "Patient ID required"
		goto renderHTML
	}

	if err := h.patients.LinkPatient(r.Context(), *currentUser, patientID); err != nil {
		if errors.Is(err, customerrors.NotFoundError) {
			comment = "Patient account for this code not found"
			goto renderHTML
		}
		if errors.Is(err, customerrors.ValidationErrorUUID) {
			comment = "Invalid code"
			goto renderHTML
		}
		comment = "Internal Error"
		goto renderHTML
	}
	comment = "Patient account created"
	success = true

renderHTML:
	w.Header().Set("HX-Trigger", "refresh-accs")
	if err := views.LinkAccountForm(comment, success).Render(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) GetPatientHandler(w http.ResponseWriter, r *http.Request) {
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
	if patient.CreatorID != currentUser.ID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	reports, err := h.rgen.GetReportsByPatient(r.Context(), patient.ID)
	if err != nil {
		http.Redirect(w, r, "/internal_server_error", http.StatusFound)
		return
	}

	if err := views.Layout(views.PatientInfo(patient, reports), patient.Name).Render(r.Context(), w); err != nil {
		http.Redirect(w, r, "/internal_server_error", http.StatusFound)
		return
	}

}

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

		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	if currentUser.IsDoctor {
		patients, err := h.patients.GetAll(r.Context(), *currentUser)
		if err != nil {
			h.logger.Error().Err(err).Msg("error getting patients")
			http.Redirect(w, r, "/internal_server_error", http.StatusFound)
			return
		}
		err = views.Layout(views.Home(patients, nil), "My Patients").Render(r.Context(), w)
		if err != nil {
			h.logger.Error().Err(err).Msg("error rendering layout")
			http.Redirect(w, r, "/internal_server_error", http.StatusFound)
		}
		return
	} else {
		reports, err := h.patients.GetApprovedReportsByUser(r.Context(), *currentUser)
		if err != nil {
			http.Redirect(w, r, "/internal_server_error", http.StatusFound)
			return
		}
		err = views.Layout(views.Home(nil, reports), "Home").Render(r.Context(), w)
		if err != nil {
			h.logger.Error().Err(err).Msg("error rendering layout")
			http.Redirect(w, r, "/internal_server_error", http.StatusFound)
		}
		return
	}
}

func (h *Handlers) GenerateReportByPatientHandler(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := GetCurrentUser(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	patientUUID, err := h.patients.GetSelfPatientIDOfUser(r.Context(), *currentUser)
	if err != nil {
		http.Redirect(w, r, "/internal_server_error", http.StatusTemporaryRedirect)
		return
	}

	file, fileHeader, err := r.FormFile("xray")
	if err != nil || fileHeader == nil {
		h.logger.Error().Err(err).Any("fileHeader", fileHeader).Msg("r.FormFile(\"xray\")")
		http.Redirect(w, r, "/internal_server_error", http.StatusFound)
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			h.logger.Error().Err(err).Msg("file.Close()")
		}
	}(file)

	_, err = h.rgen.GenerateReport(r.Context(), patientUUID, file, filepath.Ext(fileHeader.Filename))
	if err != nil {
		h.logger.Error().Err(err).Msg("rgen.GenerateReport()")
		http.Redirect(w, r, "/internal_server_error", http.StatusFound)
		return
	}

	reports, err := h.patients.ListAIGeneratedReportsByPatientID(r.Context(), patientUUID)
	if err = views.ReportsList(reports).Render(r.Context(), w); err != nil {
		w.Header().Set("HX-Redirect", "/internal_server_error")
		w.WriteHeader(http.StatusOK)
		return
	}
	return
}

func (h *Handlers) NewAiReportForPatient(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := GetCurrentUser(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	patientUUID, err := h.patients.GetSelfPatientIDOfUser(r.Context(), *currentUser)
	if err != nil {
		http.Redirect(w, r, "/internal_server_error", http.StatusTemporaryRedirect)
		return
	}

	reports, err := h.patients.ListAIGeneratedReportsByPatientID(r.Context(), patientUUID)
	if err != nil {
		http.Redirect(w, r, "/internal_server_error", http.StatusTemporaryRedirect)
		return
	}

	if err = views.Layout(views.NewAIReportPage(reports), "Generate report").Render(r.Context(), w); err != nil {
		http.Redirect(w, r, "/internal_server_error", http.StatusTemporaryRedirect)
		return
	}
}
