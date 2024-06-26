package http

import (
	"errors"
	"github.com/a-h/templ"
	"github.com/radiologist-ai/web-app/internal/app/http/handlers"
	"github.com/radiologist-ai/web-app/internal/views"
	"net/http"
)

func NewRouter(handlers *handlers.Handlers) (*http.ServeMux, error) {
	if handlers == nil {
		return nil, errors.New("handler is nil")
	}
	mux := http.NewServeMux()

	// auth
	mux.HandleFunc("GET /{$}",
		handlers.WithHTMLResponse(
			handlers.WithCurrentUser(
				handlers.HandleIndex)))
	mux.HandleFunc("GET /register",
		handlers.WithHTMLResponse(
			handlers.WithCurrentUser(
				handlers.AnonymousRequired(
					templ.Handler(
						views.Layout(
							views.RegistrationForm(),
							"Radiologist AI. Register.")).
						ServeHTTP))))
	mux.HandleFunc("POST /register",
		handlers.WithCurrentUser(
			handlers.AnonymousRequired(
				handlers.PostRegister)))
	mux.HandleFunc("POST /validate/email",
		handlers.WithHTMLResponse(
			handlers.ValidateEmail))
	mux.HandleFunc("POST /validate/password",
		handlers.WithHTMLResponse(
			handlers.ValidatePassword))

	mux.HandleFunc("GET /login",
		handlers.WithHTMLResponse(
			handlers.WithCurrentUser(
				handlers.AnonymousRequired(
					templ.Handler(
						views.Layout(
							views.LoginForm(),
							"Radiologist AI. Login.")).
						ServeHTTP))))
	mux.HandleFunc("POST /login",
		handlers.WithHTMLResponse(
			handlers.WithCurrentUser(
				handlers.AnonymousRequired(
					handlers.PostLogin))))

	mux.HandleFunc("POST /logout",
		handlers.WithCurrentUser(
			handlers.AuthRequired(
				handlers.PostLogout)))

	mux.HandleFunc("GET /home",
		handlers.WithHTMLResponse(
			handlers.WithCurrentUser(
				handlers.GetHomeHandler)))

	mux.HandleFunc("POST /my-patients",
		handlers.WithCurrentUser(
			handlers.AuthRequired(
				handlers.PostPatientHandler)))

	mux.HandleFunc("GET /patients/{patientID}/reports",
		handlers.WithHTMLResponse(
			handlers.WithCurrentUser(
				handlers.AuthRequired(
					handlers.GetPatientHandler))))

	mux.HandleFunc("GET /patients/{patientID}/reports/new",
		handlers.WithHTMLResponse(
			handlers.WithCurrentUser(
				handlers.AuthRequired(
					handlers.NewReportHandler))))

	mux.HandleFunc("POST /patients/{patientID}/reports/new",
		handlers.WithHTMLResponse(
			handlers.WithCurrentUser(
				handlers.AuthRequired(
					handlers.PostNewReportHandler))))

	mux.HandleFunc("GET /patients/{patientID}/reports/{reportID}",
		handlers.WithHTMLResponse(
			handlers.WithCurrentUser(
				handlers.AuthRequired(
					handlers.HandleEditReportPage))))
	mux.HandleFunc("POST /patients/{patientID}/reports/{reportID}/update",
		handlers.WithHTMLResponse(
			handlers.WithCurrentUser(
				handlers.AuthRequired(
					handlers.PostEditReportHandler))))

	mux.HandleFunc("POST /my-reports",
		handlers.WithHTMLResponse(
			handlers.WithCurrentUser(
				handlers.AuthRequired(
					handlers.GenerateReportByPatientHandler))))
	mux.HandleFunc("GET /ai-reports",
		handlers.WithHTMLResponse(
			handlers.WithCurrentUser(
				handlers.AuthRequired(
					handlers.NewAiReportForPatient))))

	mux.HandleFunc("GET /link-account",
		handlers.WithHTMLResponse(
			handlers.WithCurrentUser(
				handlers.AuthRequired(
					templ.Handler(
						views.Layout(
							views.LinkAccountPage(),
							"Link account")).
						ServeHTTP))))

	mux.HandleFunc("POST /link-account",
		handlers.WithHTMLResponse(
			handlers.WithCurrentUser(
				handlers.AuthRequired(
					handlers.PostLinkAccountHandler))))

	mux.HandleFunc("GET /my-accounts",
		handlers.WithHTMLResponse(
			handlers.WithCurrentUser(
				handlers.AuthRequired(
					handlers.GetMyAccountsHandler))))

	// technical
	mux.HandleFunc("GET /internal_server_error",
		handlers.WithHTMLResponse(
			templ.Handler(
				views.Layout(
					views.InternalError(),
					"Internal Error")).
				ServeHTTP))

	mux.HandleFunc("GET /", handlers.WithHTMLResponse(
		handlers.WithCurrentUser(
			templ.Handler(
				views.Layout(
					views.NotFound(),
					"404")).
				ServeHTTP)))
	return mux, nil
}
