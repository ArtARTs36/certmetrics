package internal

import (
	"log/slog"
	"time"
)

func Collect(cfg *Config) {
	collect(cfg)

	t := time.NewTicker(cfg.Inspect.Interval)
	for range t.C {
		collect(cfg)
	}
}

func collect(cfg *Config) {
	for i, pem := range cfg.Inspect.X509.PEMs {
		err := InspectX509Pem(pem)
		if err != nil {
			slog.
				With(slog.Any("pem_index", i)).
				With(slog.Any("pem_id", pem.ID)).
				With(slog.Any("err", err)).
				Warn("[collect] failed to collect x509 certs")
		}
	}
}
