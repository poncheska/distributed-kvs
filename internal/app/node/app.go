package node

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"distributed-kvs/internal/configs"
	"distributed-kvs/internal/server/http"
	storehandler "distributed-kvs/internal/server/http/handlers/store"
	"distributed-kvs/internal/store"
)

func Run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := configs.ReadLocal()
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	storeSvc, err := store.New(ctx, cfg.Store, logger)
	if err != nil {
		return fmt.Errorf("new store: %w", err)
	}

	storeHandler := storehandler.New(storeSvc, logger)

	err = http.NewServer(cfg.HTTPServer, storeHandler, logger).Start(ctx)
	if err != nil {
		return fmt.Errorf("start server: %w", err)
	}

	return nil
}
