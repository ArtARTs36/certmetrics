package scrappers

import (
	"context"

	"github.com/artarts36/certmetrics/internal/config"
)

type Scrapper interface {
	Scrape(ctx context.Context, cfg *config.ScrapeConfig) error
}
