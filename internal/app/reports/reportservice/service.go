package reportservice

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/radiologist-ai/web-app/internal/domain"
	"github.com/rs/zerolog"
	"io"
	"time"
)

type (
	ReportService struct {
		repo     domain.ReportRepository
		logger   *zerolog.Logger
		s3Client domain.S3Adapter
		rgen     domain.RGen
	}
)

func (s *ReportService) GenerateReport(ctx context.Context, patientID uuid.UUID, photo io.ReadCloser, size int64) (domain.ReportModel, error) {
	//	TODO unmock
	report, err := s.rgen.GenerateReport(ctx, "http://minio:9000/public/asdasdada.jpg")
	if err != nil {
		return domain.ReportModel{}, err
	}
	return domain.ReportModel{
		ID:         1,
		PatientID:  1,
		ImagePath:  "http://localhost:9000/public/s50010747.jpg",
		ReportText: report,
		Approved:   false,
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	}, nil
}

func New(logger *zerolog.Logger, repo domain.ReportRepository, s3Client domain.S3Adapter, rgen domain.RGen) (*ReportService, error) {
	if logger == nil {
		return nil, errors.New("logger is required")
	}
	if repo == nil {
		return nil, errors.New("repo is required")
	}
	if s3Client == nil {
		return nil, errors.New("s3Client is required")
	}
	if rgen == nil {
		return nil, errors.New("rgen is required")
	}
	return &ReportService{
		repo:     repo,
		logger:   logger,
		s3Client: s3Client,
		rgen:     rgen,
	}, nil
}
