package node

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"distributed-kvs/internal/configs"
	"distributed-kvs/internal/server/grpc"
	storegrpchandler "distributed-kvs/internal/server/grpc/handlers/store"
	"distributed-kvs/internal/server/http"
	storehttphandler "distributed-kvs/internal/server/http/handlers/store"
	"distributed-kvs/internal/store"
	"golang.org/x/sync/errgroup"
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

	storeHTTPHandler := storehttphandler.New(storeSvc, logger)

	storeGRPCHandler := storegrpchandler.New(storeSvc)

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		errG := http.NewServer(cfg.HTTPServer, storeHTTPHandler, logger).Start(ctx)
		if errG != nil {
			return fmt.Errorf("start server (%v): %w", cfg.HTTPServer.Port, errG)
		}

		return nil
	})

	eg.Go(func() error {
		errG := grpc.NewServer(cfg.GRPCServer, storeGRPCHandler, logger).Start(ctx)
		if errG != nil {
			return fmt.Errorf("start server (%v): %w", cfg.GRPCServer.Port, errG)
		}

		return nil
	})

	return nil
}
