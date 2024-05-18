package reportrepo

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/radiologist-ai/web-app/internal/domain"
	"github.com/radiologist-ai/web-app/internal/domain/customerrors"
	"github.com/rs/zerolog"
	"strings"
)

type ReportRepo struct {
	db     *sqlx.DB
	logger *zerolog.Logger
}

func NewReportRepo(db *sqlx.DB, logger *zerolog.Logger) (*ReportRepo, error) {
	if db == nil {
		return nil, errors.New("db is required")
	}
	if logger == nil {
		return nil, errors.New("logger is required")
	}
	return &ReportRepo{db, logger}, nil
}

func (r *ReportRepo) PatchReport(ctx context.Context, id int, opts ...domain.PatchOpt) error {
	if len(opts) == 0 {
		return errors.New("opts is required")
	}
	form := domain.ReportPatchForm{ID: id}
	conf := domain.PatchConf{}
	for _, opt := range opts {
		opt(&conf)
	}
	q := `UPDATE reports SET`
	qq := make([]string, 0)
	if conf.ReportText != nil {
		qq = append(qq, ` ReportText = :ReportText`)
		form.ReportText = *conf.ReportText
	}
	if conf.Approved != nil {
		qq = append(qq, ` Approved = :Approved`)
		form.Approved = *conf.Approved
	}
	q += strings.Join(qq, `,`) + ` WHERE id = :id`
	_, err := r.db.NamedExecContext(ctx, q, form)
	if err != nil {
		return fmt.Errorf("%w%w", customerrors.InternalErrorSQL, err)
	}
	return nil
}

func (r *ReportRepo) CreateReport(ctx context.Context, patientID int, imagePath, reportText string, approved bool) (createdModel domain.ReportModel, err error) {
	q := `INSERT INTO reports (patient_id, image_path, report_text, Approved) VALUES ($1, $2, $3, $4)`
	err = r.db.QueryRowxContext(ctx, q, patientID, imagePath, reportText, approved).StructScan(&createdModel)
	if err != nil {
		r.logger.Error().Err(err).Int("patientID", patientID).Str("imagePath", imagePath).Msg("QueryRowxContext")
		return createdModel, err
	}
	return createdModel, nil
}
