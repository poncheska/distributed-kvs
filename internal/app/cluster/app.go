package cluster

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"distributed-kvs/internal/configs"
	"distributed-kvs/internal/server/grpc"
	"distributed-kvs/internal/server/http"
	storehttphandler "distributed-kvs/internal/server/http/handlers/store"
	storegrpchandler "distributed-kvs/internal/server/grpc/handlers/store"
	"distributed-kvs/internal/store"
	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := configs.ReadClusterLocal()
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	eg, ctx := errgroup.WithContext(ctx)

	for i := range cfg {
		var storeSvc store.Store

		nodeCfg := cfg[i]

		nodeLogger := logger.With(slog.Int("port", nodeCfg.HTTPServer.Port))

		storeSvc, err = store.New(ctx, nodeCfg.Store, nodeLogger)
		if err != nil {
			return fmt.Errorf("new store: %w", err)
		}

		storeHTTPHandler := storehttphandler.New(storeSvc, nodeLogger)

		eg.Go(func() error {
			errG := http.NewServer(nodeCfg.HTTPServer, storeHTTPHandler, nodeLogger).Start(ctx)
			if errG != nil {
				return fmt.Errorf("start server (%v): %w", nodeCfg.HTTPServer.Port, errG)
			}

			return nil
		})

		storeGRPCHandler := storegrpchandler.New(storeSvc)

		eg.Go(func() error {
			errG := grpc.NewServer(nodeCfg.GRPCServer, storeGRPCHandler, nodeLogger).Start(ctx)
			if errG != nil {
				return fmt.Errorf("start server (%v): %w", nodeCfg.GRPCServer.Port, errG)
			}

			return nil
		})
	}

	err = eg.Wait()
	if err != nil {
		return fmt.Errorf("wait: %w", err)
	}

	return nil
}
