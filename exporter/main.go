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

	"github.com/artarts36/certmetrics/exporter/internal"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	httpReadTimeout = 3 * time.Second
	shutdownTimeout = 5 * time.Second
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	slog.Info("loading config")
	cfg, err := internal.LoadConfig("certmetrics.yaml")
	if err != nil {
		slog.
			With(slog.Any("err", err)).
			Error("failed to load config")
		os.Exit(1)
	}

	slog.Info("setup collector")

	if err = internal.SetupCollector(); err != nil {
		slog.
			With(slog.Any("err", err)).
			Error("failed to setup collector")
		os.Exit(1)
	}

	go internal.Collect(cfg)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:              cfg.HTTP.Addr,
		Handler:           mux,
		ReadHeaderTimeout: httpReadTimeout,
	}
	go shutdown(srv)

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

func shutdown(s *http.Server) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(ch)

	sig := <-ch
	slog.
		With(slog.String("signal", sig.String())).
		Info("shutdown..")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		slog.
			With(slog.Any("err", err.Error())).
			Error("failed to shutdown http server")
	}
}
