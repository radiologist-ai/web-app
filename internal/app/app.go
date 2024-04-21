package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	http2 "github.com/radiologist-ai/web-app/internal/app/http"
	"github.com/radiologist-ai/web-app/internal/app/http/handlers"
	"github.com/radiologist-ai/web-app/internal/app/users/usersrepo"
	"github.com/radiologist-ai/web-app/internal/app/users/usersservice"
	"github.com/radiologist-ai/web-app/internal/config"
	"github.com/radiologist-ai/web-app/pkg/ptr"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"sync"
	"time"
)

func Run(backgroundCtx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()

	cfg := config.GetConfig()

	logger := ptr.Pointer(zerolog.New(os.Stderr).With().Timestamp().Caller().Logger())
	// TODO Unmock
	db := &sqlx.DB{
		DB:     nil,
		Mapper: nil,
	}

	// repository
	usersRepo, err := usersrepo.New(logger, db)
	if err != nil {
		return err
	}

	// service
	usersService, err := usersservice.New(logger, usersRepo)
	if err != nil {
		return err
	}

	// handlers
	handle, err := handlers.NewHandlers(logger, usersService)
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

//func PreparePostgres(cfg)