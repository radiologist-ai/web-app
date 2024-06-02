package patientrepo

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/radiologist-ai/web-app/internal/domain"
	"github.com/radiologist-ai/web-app/internal/domain/customerrors"
)

func (pr *PatientRepo) SelectAIGeneratedReportsByPatient(ctx context.Context, patientID uuid.UUID) ([]domain.ReportForPatient, error) {
	res := make([]domain.ReportForPatient, 0)
	query := `
		SELECT 
		    r.id as id, 
		    p.name as name, 
		    r.image_path as image_path, 
		    r.report_text as report_text,
			r.approved as approved,
			r.created_at as created_at,
			r.updated_at as updated_at
		FROM reports AS r INNER JOIN 
		    patients AS p ON p.id = r.patient_id
		WHERE p.id = $1
		ORDER BY r.updated_at DESC
`
	err := pr.db.SelectContext(ctx, &res, query, patientID)
	if err != nil {
		return nil, fmt.Errorf("%w%w", customerrors.InternalErrorSQL, err)
	}
	return res, nil
}

func (pr *PatientRepo) SelectApprovedReportsByUser(ctx context.Context, userID int) ([]domain.ReportForPatient, error) {
	res := make([]domain.ReportForPatient, 0)
	query := `
		SELECT 
		    r.id as id, 
		    p.name as name,
		    u.first_name || ' ' || u.last_name as doctor_full_name,
		    r.image_path as image_path, 
		    r.report_text as report_text,
			r.approved as approved,
			r.created_at as created_at,
			r.updated_at as updated_at
		FROM reports AS r INNER JOIN 
		    patients AS p ON p.id = r.patient_id INNER JOIN 
		    users u on u.id = p.creator_id
		WHERE p.user_id = $1 AND p.creator_id != $1 AND r.approved = true
		ORDER BY r.updated_at DESC
`
	err := pr.db.SelectContext(ctx, &res, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%w%w", customerrors.InternalErrorSQL, err)
	}
	return res, nil
}
