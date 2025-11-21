package internal

import (
	"context"
	"fmt"
	"github.com/artarts36/certmetrics/pkg/collector"
	"log/slog"
	"regexp"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/artarts36/certmetrics/internal/config"
	"github.com/artarts36/certmetrics/internal/metrics"
	"github.com/artarts36/certmetrics/internal/scrappers"
	"github.com/artarts36/certmetrics/internal/storage"
)

type App struct {
	cfg      *config.Config
	scrapper scrappers.Scrapper

	exporterMetrics *metrics.ExporterMetrics
}

type AppInfo struct {
	Version   string
	BuildDate string
}

var (
	regexDomain = regexp.MustCompile(`^https:\/\/[a-zA-Z0-9-.]+$`)
	regexHTTP   = regexp.MustCompile(`^http:\/\/`)
	regexHTTPS  = regexp.MustCompile(`^https:\/\/(.*)/(.*)`)
	regexENV    = regexp.MustCompile(`^env:\/\/`)
)

func NewApp(cfg *config.Config, info AppInfo) (*App, error) {
	app := &App{
		cfg:             cfg,
		exporterMetrics: metrics.NewExporterMetrics("certmetrics_exporter"),
	}

	store := storage.NewResolver(
		storage.NewLocal(),
		[]*storage.ResolveRule{
			{
				Regex:   regexDomain,
				Storage: storage.NewDomain(),
			},
			{
				Regex:   regexHTTP,
				Storage: storage.NewHTTP(),
			},
			{
				Regex:   regexHTTPS,
				Storage: storage.NewHTTP(),
			},
			{
				Regex:   regexENV,
				Storage: storage.NewEnv(),
			},
		},
	)

	sc := map[string]scrappers.Scrapper{}

	if len(cfg.Scrape.X509.Files) > 0 {
		sc["x509"] = scrappers.NewX509Scrapper(app.exporterMetrics, store)
	}

	app.scrapper = scrappers.Parallel(sc)

	if err := app.setupMetrics(); err != nil {
		return nil, fmt.Errorf("setupMetrics: %w", err)
	}

	app.exporterMetrics.SetInfo(info.Version, info.BuildDate)

	return app, nil
}

func (app *App) Run(ctx context.Context) {
	app.scrape(ctx)

	t := time.NewTicker(app.cfg.Scrape.Interval.Value)
	for range t.C {
		app.scrape(ctx)
	}
}

func (app *App) setupMetrics() error {
	collector.DefaultCollector.As("exporter")
	collector.DefaultCollector.WithExpiryDays()

	if err := collector.Register(); err != nil {
		return fmt.Errorf("register metrics: %w", err)
	}
	if err := prometheus.Register(app.exporterMetrics); err != nil {
		return fmt.Errorf("register exporter metrics: %w", err)
	}

	return nil
}

func (app *App) scrape(ctx context.Context) {
	slog.DebugContext(ctx, "[app] running scrappers")

	err := app.scrapper.Scrape(ctx, &app.cfg.Scrape)
	if err != nil {
		slog.
			With(slog.Any("err", err)).
			ErrorContext(ctx, "failed to scrap")
	}
}
