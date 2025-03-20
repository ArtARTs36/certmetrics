package internal

import (
	"fmt"

	"github.com/artarts36/certmetrics"
)

func SetupCollector() error {
	certmetrics.DefaultCollector.As("exporter")

	if err := certmetrics.Register(); err != nil {
		return fmt.Errorf("register metrics: %w", err)
	}

	return nil
}
