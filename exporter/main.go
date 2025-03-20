package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/artarts36/certmetrics"
	"github.com/artarts36/certmetrics/x509m"
)

const (
	httpReadTimeout = 3 * time.Second
	shutdownTimeout = 5 * time.Second
)

type Config struct {
	HTTP struct {
		Addr string `yaml:"addr"`
	} `yaml:"http"`

	Inspect struct {
		Interval time.Duration `yaml:"interval"`

		X509 struct {
			// Paths to .pem
			PEMs []PEMFile `yaml:"pems"`
		} `yaml:"x509"`
	} `yaml:"inspect"`
}

type PEMFile struct {
	Path string `yaml:"path"`
	ID   string `yaml:"id"`
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	slog.Info("loading config")
	cfg, err := loadConfig("certmetrics.yaml")
	if err != nil {
		slog.
			With(slog.Any("err", err)).
			Error("failed to load config")
		os.Exit(1)
	}

	slog.Info("registering metrics")

	if err = certmetrics.Register(); err != nil {
		slog.
			With(slog.Any("err", err)).
			Error("failed to register metrics")
		os.Exit(1)
	}

	go func() {
		collect(cfg)

		t := time.NewTicker(cfg.Inspect.Interval)
		for range t.C {
			collect(cfg)
		}
	}()

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

func collect(cfg *Config) {
	for i, pem := range cfg.Inspect.X509.PEMs {
		err := inspectPem(pem)
		if err != nil {
			slog.
				With(slog.Any("pem_id", i)).
				With(slog.Any("err", err)).
				Warn("[collect] failed to collect x509 certs")
		}
	}
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

func inspectPem(pem PEMFile) error {
	file, err := os.ReadFile(pem.Path)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	x509m.InspectPEMs(file, x509m.WithID(pem.ID))

	return nil
}

func loadConfig(path string) (*Config, error) {
	var cfg Config

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	if err = yaml.Unmarshal(content, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal yaml: %w", err)
	}

	return &cfg, nil
}
