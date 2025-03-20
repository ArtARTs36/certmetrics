package scrappers

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/artarts36/certmetrics/exporter/internal/config"
)

type parallel struct {
	scrappers map[string]Scrapper
}

func Parallel(scrappers map[string]Scrapper) Scrapper {
	return &parallel{
		scrappers: scrappers,
	}
}

func (c *parallel) Scrape(ctx context.Context, cfg *config.Config) error {
	wg := &sync.WaitGroup{}
	gerr := []error{}

	for name, scrapper := range c.scrappers {
		wg.Add(1)

		go func(name string, sc Scrapper) {
			defer wg.Done()

			err := sc.Scrape(ctx, cfg)
			if err != nil {
				gerr = append(gerr, fmt.Errorf("%s scrapper: %w", name, err))
			}
		}(name, scrapper)
	}

	return errors.Join(gerr...)
}
