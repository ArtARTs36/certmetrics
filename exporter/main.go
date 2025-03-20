package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/artarts36/certmetrics/exporter/internal"
	"github.com/artarts36/certmetrics/exporter/internal/config"
)

var (
	Version   = "v0.1.0"
	BuildDate = ""
)

const (
	httpReadTimeout = 3 * time.Second
	shutdownTimeout = 5 * time.Second
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	exit := func() {
		cancel()
		os.Exit(1)
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	slog.Info("loading config")
	cfg, err := config.Load("certmetrics.yaml")
	if err != nil {
		slog.
			With(slog.Any("err", err)).
			Error("failed to load config")
		exit()
	}

	slog.Info("initializing application")

	app, err := internal.NewApp(cfg, internal.AppInfo{
		Version:   Version,
		BuildDate: BuildDate,
	})
	if err != nil {
		slog.
			With(slog.Any("err", err)).
			Error("failed to load config")
		exit()
	}

	go app.Run(ctx)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:              cfg.HTTP.Addr,
		Handler:           mux,
		ReadHeaderTimeout: httpReadTimeout,
	}
	go shutdown(srv, cancel)

	slog.
		With(slog.String("addr", cfg.HTTP.Addr)).
		Info("listening http")

	if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.
			With(slog.Any("err", err.Error())).
			Error("failed to listen http")
	}

	slog.Info("http server closed")
}

func shutdown(s *http.Server, cancel context.CancelFunc) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(ch)

	sig := <-ch
	slog.
		With(slog.String("signal", sig.String())).
		Info("shutdown..")

	ctx, shCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shCancel()

	if err := s.Shutdown(ctx); err != nil {
		slog.
			With(slog.Any("err", err.Error())).
			Error("failed to shutdown http server")
	}

	cancel()
}
