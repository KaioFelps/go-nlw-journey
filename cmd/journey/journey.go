package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/phenpessoa/gutils/netutils/httputils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"nlw-journey/internal/api"
	"nlw-journey/internal/api/spec"
	"nlw-journey/internal/mail/mailpit"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill, syscall.SIGKILL, syscall.SIGTERM)
	defer cancel()

	if err := godotenv.Load(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := run(ctx); err != nil {
		fmt.Println("Server forcedly closed because of an error:")
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Println("Shutting down the application.")
}

func run(ctx context.Context) error {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, err := cfg.Build()
	if err != nil {
		return err
	}

	logger = logger.Named("journey_api")
	defer func() {
		_ = logger.Sync()
	}()

	pool, err := pgxpool.New(ctx, fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s",
		os.Getenv("JOURNEY_DATABASE_HOST"),
		os.Getenv("JOURNEY_DATABASE_PORT"),
		os.Getenv("JOURNEY_DATABASE_USER"),
		os.Getenv("JOURNEY_DATABASE_NAME"),
		os.Getenv("JOURNEY_DATABASE_PASSWORD"),
	),
	)

	if err != nil {
		return err
	}

	if err := pool.Ping(ctx); err != nil {
		return err
	}

	var mailer = mailpit.NewMailPit(pool, logger)

	si := api.NewAPI(pool, logger, mailer)

	router := chi.NewRouter()
	router.Use(middleware.RequestID, middleware.Recoverer, httputils.ChiLogger(logger))
	router.Mount("/", spec.Handler(si))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	defer func() {
		const timeout = 30 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Error("Failed to shutdown server", zap.Error(err))
		}
	}()

	errChan := make(chan error, 1)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			errChan <- err
		}
	}()

	go func() {
		time.Sleep(time.Second)
		logger.Info("Server running on port 8080.")
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errChan:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}

	return nil
}
