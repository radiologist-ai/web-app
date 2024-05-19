package reportservice

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/radiologist-ai/web-app/internal/domain"
	"github.com/radiologist-ai/web-app/internal/domain/customerrors"
	"github.com/radiologist-ai/web-app/pkg/image_processing"
	"github.com/radiologist-ai/web-app/pkg/s3ObjectNameGenerator"
	"github.com/rs/zerolog"
	"image"
	"io"
)

type (
	ReportService struct {
		repo     domain.ReportRepository
		logger   *zerolog.Logger
		s3Client domain.S3Adapter
		rgen     domain.RGen
	}
)

func (s *ReportService) GenerateReport(ctx context.Context, patientID uuid.UUID, photo io.Reader, ext string) (domain.ReportModel, error) {
	var (
		resCh = make(chan string)
		errCh = make(chan error)
	)
	if ext == "" {
		return domain.ReportModel{}, fmt.Errorf("%winvalid filename without extention", customerrors.ValidationError)
	}
	img, sniff, err := image.Decode(photo)
	if err != nil {
		return domain.ReportModel{}, err
	}

	toSave, err := image_processing.ProcessImage(img, domain.XrayWidth, domain.XrayHeight, sniff)
	if err != nil {
		return domain.ReportModel{}, err
	}

	objName := s3ObjectNameGenerator.NewObjectName(fmt.Sprintf("patients/%s/reports", patientID.String()), ext)
	realObjName, err := s.s3Client.Save(ctx, objName, "application/octet-stream", toSave, toSave.Size())
	if err != nil {
		return domain.ReportModel{}, err
	}

	go s.rgen.GenerateReportAsync(ctx, s.s3Client.GetPublicLink(realObjName), resCh, errCh)

	select {
	case err := <-errCh:
		return domain.ReportModel{}, err
	case report := <-resCh:
		result, err2 := s.repo.CreateReport(ctx, patientID, realObjName, report, false)
		if err2 != nil {
			return domain.ReportModel{}, err2
		}
		return result, nil
	case <-ctx.Done():
		return domain.ReportModel{}, ctx.Err()
	}
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
