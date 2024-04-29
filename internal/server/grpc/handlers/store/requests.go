package store

import (
	"context"

	"distributed-kvs/pkg/api/store"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) Get(ctx context.Context, request *store.GetRequest) (*store.GetResponse, error) {
	val, err := h.store.Get(request.GetKey())
	if err != nil {

	}

	return &store.GetResponse{Value: val}, nil
}

func (h *Handler) Set(ctx context.Context, request *store.SetRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (h *Handler) Delete(ctx context.Context, request *store.DeleteRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (h *Handler) Join(ctx context.Context, request *store.JoinRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
