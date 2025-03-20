package scrappers

import (
	"context"
	"fmt"
	"os"

	"github.com/artarts36/certmetrics/exporter/internal/config"
	"github.com/artarts36/certmetrics/exporter/internal/metrics"
	"github.com/artarts36/certmetrics/x509m"
)

type X509Scrapper struct {
	metrics *metrics.ExporterMetrics
}

func NewX509Scrapper(metr *metrics.ExporterMetrics) *X509Scrapper {
	return &X509Scrapper{
		metrics: metr,
	}
}

func (x *X509Scrapper) Scrape(_ context.Context, cfg *config.Config) error {
	for i, pem := range cfg.Scrape.X509.PEMs {
		err := x.scrape(pem)
		if err != nil {
			return fmt.Errorf("scrape file with index %d: %w", i, err)
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
