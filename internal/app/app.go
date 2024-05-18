package app

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	http2 "github.com/radiologist-ai/web-app/internal/app/http"
	"github.com/radiologist-ai/web-app/internal/app/http/handlers"
	"github.com/radiologist-ai/web-app/internal/app/patient/patientrepo"
	"github.com/radiologist-ai/web-app/internal/app/patient/patientservice"
	"github.com/radiologist-ai/web-app/internal/app/reports/reportrepo"
	"github.com/radiologist-ai/web-app/internal/app/reports/reportservice"
	"github.com/radiologist-ai/web-app/internal/app/users/usersrepo"
	"github.com/radiologist-ai/web-app/internal/app/users/usersservice"
	"github.com/radiologist-ai/web-app/internal/config"
	"github.com/radiologist-ai/web-app/pkg/minio_adapter"
	"github.com/radiologist-ai/web-app/pkg/ptr"
	"github.com/radiologist-ai/web-app/pkg/rgen"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"os"
	"sync"
	"time"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func Run(backgroundCtx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()

	cfg := config.GetConfig()

	logger := ptr.Pointer(zerolog.New(os.Stderr).With().Timestamp().Caller().Logger())

	db, err := PreparePostgres(cfg.Database)
	if err != nil {
		return err
	}

	grpcConn, err := grpc.DialContext(backgroundCtx, cfg.GRPC.Addr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	//third-party shenanigans
	rgenCli, err := rgen.NewClient(grpcConn)
	if err != nil {
		return err
	}
	mc, err := minio_adapter.New(logger, &minio_adapter.MinioConfig{
		ServerURL:  cfg.Minio.ServerURL,
		AccessKey:  cfg.Minio.AccessKey,
		SecretKey:  cfg.Minio.SecretKey,
		BucketName: cfg.Minio.BucketName,
	}, false)
	if err != nil {
		return err
	}

	// repository
	usersRepo, err := usersrepo.New(logger, db)
	if err != nil {
		return err
	}
	patientRepo, err := patientrepo.New(db, logger)
	if err != nil {
		return err
	}
	reportRepo, err := reportrepo.NewReportRepo(db, logger)
	if err != nil {
		return err
	}

	// service
	usersService, err := usersservice.New(logger, usersRepo)
	if err != nil {
		return err
	}
	patientService, err := patientservice.New(logger, patientRepo)
	if err != nil {
		return err
	}
	reportService, err := reportservice.New(logger, reportRepo, mc, rgenCli)

	// handlers
	handle, err := handlers.NewHandlers(logger, usersService, patientService, reportService, cfg.Server.Secret)
	if err != nil {
		return err
	}
	router, err := http2.NewRouter(handle)
	if err != nil {
		return err
	}

	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Server.ListenAddr, cfg.Server.Port),
		Handler: router,
	}

	errChan := make(chan error)
	logger.Info().Msg(fmt.Sprintf("starting server on %s:%d", cfg.Server.ListenAddr, cfg.Server.Port))
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	select {
	case <-backgroundCtx.Done():
		fmt.Println("Shutting down HTTP server gracefully...")
		shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelShutdown()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			return err
		}
		return nil
	case err := <-errChan:
		if err != nil {
			return err
		}
	}

	logger.Info().Msg("HTTP server stopped.")
	return nil
}

func PreparePostgres(cfg config.Database) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres",
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database))
	if err != nil {
		return nil, err
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return nil, err
	}

	if err := goose.Up(db.DB, "migrations"); err != nil {
		return nil, err
	}
	return db, nil
}
