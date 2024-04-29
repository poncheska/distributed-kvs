package store

import (
	"distributed-kvs/pkg/api/store"
)

type Store interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Delete(key string) error
	Join(nodeID string, addr string) error
}

type Handler struct {
	store Store

	store.UnimplementedStoreServer
}

func New(store Store) *Handler {
	return &Handler{
		store: store,
	}
}
