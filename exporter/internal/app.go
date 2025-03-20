package internal

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/artarts36/certmetrics"
	"github.com/artarts36/certmetrics/exporter/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/artarts36/certmetrics/exporter/internal/config"
	"github.com/artarts36/certmetrics/exporter/internal/scrappers"
)

type App struct {
	cfg      *config.Config
	scrapper scrappers.Scrapper

	exporterMetrics *metrics.ExporterMetrics
}

func NewApp(cfg *config.Config) (*App, error) {
	app := &App{
		cfg:             cfg,
		exporterMetrics: metrics.NewExporterMetrics("certmetrics_exporter"),
	}

	sc := map[string]scrappers.Scrapper{}

	if len(cfg.Scrape.X509.PEMs) > 0 {
		sc["x509"] = scrappers.NewX509Scrapper(app.exporterMetrics)
	}

	app.scrapper = scrappers.Parallel(sc)

	if err := app.setup(); err != nil {
		return nil, fmt.Errorf("setup: %w", err)
	}

	return app, nil
}

func (app *App) Run(ctx context.Context) {
	app.scrape(ctx)

	t := time.NewTicker(app.cfg.Scrape.Interval)
	for range t.C {
		app.scrape(ctx)
	}
}

func (app *App) setup() error {
	certmetrics.DefaultCollector.As("exporter")

	if err := certmetrics.Register(); err != nil {
		return fmt.Errorf("register metrics: %w", err)
	}
	if err := prometheus.Register(app.exporterMetrics); err != nil {
		return fmt.Errorf("register exporter metrics: %w", err)
	}

	return nil
}

func (app *App) scrape(ctx context.Context) {
	err := app.scrapper.Scrape(ctx, app.cfg)
	if err != nil {
		slog.
			With(slog.Any("err", err)).
			ErrorContext(ctx, "failed to scrap")
	}
}
