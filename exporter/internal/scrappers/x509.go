package scrappers

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/artarts36/certmetrics/exporter/internal/storage"

	"github.com/artarts36/certmetrics/exporter/internal/config"
	"github.com/artarts36/certmetrics/exporter/internal/metrics"
	"github.com/artarts36/certmetrics/x509m"
)

type X509Scrapper struct {
	metrics         *metrics.ExporterMetrics
	storageResolver *storage.Resolver
}

type x509Task struct {
	file  config.X509File
	store storage.Storage
}

func NewX509Scrapper(metr *metrics.ExporterMetrics, store *storage.Resolver) *X509Scrapper {
	return &X509Scrapper{
		metrics:         metr,
		storageResolver: store,
	}
}

func (x *X509Scrapper) Scrape(ctx context.Context, cfg *config.ScrapeConfig) error {
	queue := make(chan x509Task)
	wg := &sync.WaitGroup{}

	go func() {
		for task := range queue {
			err := x.scrape(ctx, task.store, task.file)
			if err != nil {
				slog.
					With(slog.Any("err", err)).
					With(slog.String("file.id", task.file.ID)).
					With(slog.String("file.path", task.file.Path)).
					ErrorContext(ctx, "[x509] failed to scrape file")
			}
			wg.Done()
		}
	}()

	for _, pem := range cfg.X509.Paths {
		store := x.storageResolver.Resolve(pem.Path)

		files, err := store.ListFiles(ctx, pem.Path)
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
			queue <- x509Task{
				file: config.X509File{
					Path: file,
					ID:   pem.ID,
					Opts: pem.Opts,
				},
				store: store,
			}
		}
	}

	close(queue)

	wg.Wait()

	return nil
}

func (x *X509Scrapper) scrape(ctx context.Context, store storage.Storage, pem config.X509File) error {
	file, err := store.ReadFile(ctx, pem.Path)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	id := pem.ID
	if id == "" {
		id = pem.Path
	}

	opts := []x509m.InspectOption{
		x509m.WithID(id),
	}

	if pem.Opts.Subject == config.PemSubjectNameOptNone {
		opts = append(opts, x509m.WithoutSubjectName())
	}

	if err = x509m.InspectPEMs(file, opts...); err != nil {
		if terr := x509m.Inspect(file, opts...); terr != nil {
			return fmt.Errorf("inspect: %w, inspest pems: %w", err, terr)
		}
		err = nil
	}

	x.metrics.IncScrapings(id)

	return nil
}
