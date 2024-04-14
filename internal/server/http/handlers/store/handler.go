package store

import (
	"errors"
	"log/slog"
	"net/http"
)

const storageKeyPathKey = "key"

type Store interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Delete(key string) error
	Join(nodeID string, addr string) error
}

type Implementation struct {
	store Store
	log   *slog.Logger
}

func New(store Store, log *slog.Logger) *Implementation {
	return &Implementation{
		store: store,
		log:   log,
	}
}

func getKeyFromPath(r *http.Request) (string, error) {
	key := r.PathValue(storageKeyPathKey)

	if key == "" {
		return "", errors.New("key must not be empty")
	}

	return key, nil
}
