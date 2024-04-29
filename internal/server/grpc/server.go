package grpc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"

	"distributed-kvs/internal/configs"
	"distributed-kvs/internal/server/grpc/interceptors"
	storeerrors "distributed-kvs/internal/store/errors"
	"distributed-kvs/pkg/api/store"
	grpclogging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	server *grpc.Server
	port   int

	log *slog.Logger
}

func NewServer(cfg configs.GRPCServerConfig, storeHandler store.StoreServer, log *slog.Logger) *Server {
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcprometheus.UnaryServerInterceptor,
			grpclogging.UnaryServerInterceptor(interceptors.InterceptorLogger(log)),
		),
	)

	store.RegisterStoreServer(s, storeHandler)

	return &Server{
		server: s,
		port:   cfg.Port,
		log:    log,
	}
}

func (s *Server) Start(ctx context.Context) error {
	var idleConnsClosed = make(chan struct{})

	go func() {
		<-ctx.Done()

		s.log.Info("grpc server shutdown")

		s.server.GracefulStop()

		close(idleConnsClosed)
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("net listen: %w", err)
	}

	s.log.Info("grpc server started")

	if err = s.server.Serve(lis); err != nil {
		return fmt.Errorf("grpc server serve: %w", err)
	}

	<-idleConnsClosed

	return nil
}

func wrapGRPCError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, storeerrors.ErrActionUnavailable):
		return status.Error(codes.Unimplemented, err.Error())

	default:
		return status.Error(codes.Internal, err.Error())
	}
}
