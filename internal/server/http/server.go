package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"distributed-kvs/docs"
	"distributed-kvs/internal/configs"
	"distributed-kvs/internal/server/http/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

var (
	readHeaderTimeout = 5 * time.Second
	idleTimeout       = 30 * time.Second

	shutdownTimeout = 30 * time.Second
)

type StoreHandler interface {
	Get(http.ResponseWriter, *http.Request)
	Set(http.ResponseWriter, *http.Request)
	Delete(http.ResponseWriter, *http.Request)
	Join(http.ResponseWriter, *http.Request)
}

// Server godoc
// @title		Distributed Store
// @version		1.0
// @description	distributed key-value storage
// @host		127.0.0.1:8080
// @schemes		http
type Server struct {
	httpServer *http.Server
	log        *slog.Logger
}

func NewServer(cfg configs.HTTPServerConfig, storeHandler StoreHandler, log *slog.Logger) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /store/{key}", storeHandler.Get)
	mux.HandleFunc("PUT /store/{key}", storeHandler.Set)
	mux.HandleFunc("DELETE /store/{key}", storeHandler.Delete)
	mux.HandleFunc("POST /store", storeHandler.Join)

	//docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%v", cfg.Port)

	mux.HandleFunc("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://%s/swagger/doc.json", docs.SwaggerInfo.Host)),
	))

	var httpHandler http.Handler

	httpHandler = middleware.WithCORSMiddleware(mux)

	if cfg.Logging {
		httpHandler = middleware.WithLoggingMiddleware(httpHandler, log)
	}

	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%v", cfg.Port),
		Handler:           httpHandler,
		ReadHeaderTimeout: readHeaderTimeout,
		IdleTimeout:       idleTimeout,
	}

	return &Server{
		httpServer: httpServer,
		log:        log,
	}
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		<-ctx.Done()

		shutdownCtx, done := context.WithTimeout(context.Background(), shutdownTimeout)
		defer done()

		s.log.Info("http server shutdown")

		if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
			s.log.Warn("http server shutdown: ", err)
		}
	}()

	s.log.Info("http server started")

	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http server listen and serve: %w", err)
	}

	return nil
}
