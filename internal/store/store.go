package store

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"distributed-kvs/internal/configs"
	"distributed-kvs/internal/store/raft"
	"github.com/google/uuid"
)

var (
	ErrInvalidStoreType = errors.New("invalid store type")
	ErrEmptyConfig      = errors.New("empty config")
)

type Store interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Delete(key string) error
	Join(nodeID string, addr string) error
}

func New(ctx context.Context, cfg configs.StoreConfig, logger *slog.Logger) (Store, error) {
	switch cfg.Type {
	case configs.RaftStoreType:
		if cfg.Raft == nil {
			return nil, fmt.Errorf("%w for %s store type", ErrEmptyConfig, cfg.Type)
		}

		store := raft.New(cfg.Raft.InMem, logger)

		store.RaftBind = cfg.Raft.Addr

		nodeID := cfg.Raft.NodeID

		if nodeID == "" {
			nodeID = uuid.New().String()
		}

		err := store.Open(cfg.Raft.EnableSingle, nodeID)
		if err != nil {
			return nil, fmt.Errorf("open raft store: %w", err)
		}

		joinRequestAsyncWithRetry(ctx, cfg.JoinURL, nodeID, store.RaftBind, logger)

		return store, nil
	case configs.ZabStoreType:
		if cfg.Zab == nil {
			return nil, fmt.Errorf("%w for %s store type", ErrEmptyConfig, cfg.Type)
		}

		return nil, nil // nolint // TODO
	default:
		return nil, ErrInvalidStoreType
	}
}

func joinRequestAsyncWithRetry(ctx context.Context, url, nodeID, addr string, logger *slog.Logger) {
	go func() {
		timeoutTimer := time.NewTimer(time.Minute)
		defer timeoutTimer.Stop()

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		select {
		case <-ticker.C:
			err := joinRequest(url, nodeID, addr)
			if err == nil {
				return
			}

			logger.Error(fmt.Sprintf("join request: %s", err.Error()))
		case <-timeoutTimer.C:
			return
		case <-ctx.Done():
			return
		}
	}()
}

func joinRequest(url, nodeID, addr string) error {
	bs, err := json.Marshal(map[string]string{
		"node_id": nodeID,
		"addr":    addr,
	})
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(bs)) //nolint:gosec // ok
	if err != nil {
		return fmt.Errorf("post request: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response status code: %v", resp.StatusCode)
	}

	return nil
}
