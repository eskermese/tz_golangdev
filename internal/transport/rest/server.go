package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/eskermese/tz_golangdev/internal/config"
	restHandler "github.com/eskermese/tz_golangdev/internal/transport/rest/handlers"
	"github.com/eskermese/tz_golangdev/pkg/logger"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
)

type Server struct {
	CertFile, KeyFile *string
	httpServer        *fiber.App
	Addr              string
	Handler           *restHandler.Handler
	idleConnsClosed   chan struct{}
	logger            logger.Logger
}

func NewServer(cfg *config.Config, handler *restHandler.Handler, logger logger.Logger) *Server {
	s := &Server{
		idleConnsClosed: make(chan struct{}),
		httpServer: fiber.New(fiber.Config{
			AppName:      "tz_golangdev",
			ReadTimeout:  cfg.ClientHTTP.ReadTimeout,
			WriteTimeout: cfg.ClientHTTP.WriteTimeout,
			BodyLimit:    cfg.ClientHTTP.MaxHeaderMegabytes << 20,
			JSONDecoder: func(data []byte, v interface{}) error {
				return jsoniter.Unmarshal(data, v)
			},
			JSONEncoder: func(v interface{}) ([]byte, error) {
				return jsoniter.Marshal(v)
			},
		}),
		Handler: handler,
		Addr:    fmt.Sprintf(":%d", cfg.ClientHTTP.Port),
		logger:  logger,
	}

	handler.InitRouter(cfg, s.httpServer)

	return s
}

func (s *Server) IsInsecure() bool {
	return s.CertFile == nil && s.KeyFile == nil
}

func (s *Server) Run(ctx context.Context) error {
	s.logger.Info(fmt.Sprintf(`serving ClientHTTP on "%s"`, s.Addr))

	go s.gracefulShutdown(ctx, s.httpServer)

	var err error
	if s.IsInsecure() {
		err = s.httpServer.Listen(s.Addr)
	} else {
		err = s.httpServer.ListenTLS(s.Addr, *s.CertFile, *s.KeyFile)
	}

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	s.Wait()

	return nil
}

func (s *Server) gracefulShutdown(ctx context.Context, httpSrv *fiber.App) {
	defer close(s.idleConnsClosed)
	<-ctx.Done()

	s.logger.Info("shutting down HTTP server")

	if err := httpSrv.Shutdown(); err != nil {
		s.logger.Error("shutting down HTTP server", logger.Error(err))
	}
}

func (s *Server) Wait() {
	<-s.idleConnsClosed
	s.logger.Info("HTTP server has processed all idle connections")
}
