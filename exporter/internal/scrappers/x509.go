package scrappers

import (
	"context"
	"fmt"
	"os"

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
	for _, pem := range cfg.X509.PEMs {
		files, err := x.storage.ListFiles(ctx, pem.Path)
		if err != nil {
			return fmt.Errorf("list files in %q: %w", pem.Path, err)
		}

		if len(files) == 0 {
			return fmt.Errorf("files not found in %q", pem.Path)
		}

		for _, file := range files {
			fpem := config.PEMFile{
				Path: file,
				ID:   pem.ID,
			}

			err = x.scrape(fpem)
			if err != nil {
				return fmt.Errorf("scrape file %q: %w", fpem.Path, err)
			}
		}
	}

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
