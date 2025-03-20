package scrappers

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/artarts36/certmetrics/exporter/internal/storage"

	"github.com/artarts36/certmetrics/exporter/internal/config"
	"github.com/artarts36/certmetrics/exporter/internal/metrics"
	"github.com/artarts36/certmetrics/x509m"
)

type X509Scrapper struct {
	metrics *metrics.ExporterMetrics
	storage storage.Storage
}

func NewX509Scrapper(metr *metrics.ExporterMetrics, store storage.Storage) *X509Scrapper {
	return &X509Scrapper{
		metrics: metr,
		storage: store,
	}
}

func (x *X509Scrapper) Scrape(ctx context.Context, cfg *config.ScrapeConfig) error {
	queue := make(chan config.PEMFile)
	wg := &sync.WaitGroup{}

	go func() {
		for file := range queue {
			err := x.scrape(file)
			if err != nil {
				slog.
					With(slog.String("file.id", file.ID)).
					With(slog.String("file.path", file.Path)).
					ErrorContext(ctx, "[x509] failed to scrape file")
			}
			wg.Done()
		}
	}()

	for _, pem := range cfg.X509.PEMs {
		files, err := x.storage.ListFiles(ctx, pem.Path)
		if err != nil {
			slog.
				With(slog.Any("err", err)).
				With(slog.String("path", pem.Path)).
				ErrorContext(ctx, "[x509] failed to list files")
			continue
		}

		if len(files) == 0 {
			slog.
				With(slog.Any("err", err)).
				With(slog.String("path", pem.Path)).
				WarnContext(ctx, "[x509] files not found")
			continue
		}

		wg.Add(len(files))

		for _, file := range files {
			queue <- config.PEMFile{
				Path: file,
				ID:   pem.ID,
			}
		}
	}

	close(queue)

	wg.Wait()

	return nil
}

func (x *X509Scrapper) scrape(pem config.PEMFile) error {
	file, err := os.ReadFile(pem.Path)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	id := pem.ID
	if id == "" {
		id = pem.Path
	}

	x509m.InspectPEMs(file, x509m.WithID(id))
	x.metrics.IncScrapings(id)

	return nil
}
