package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/eskermese/tz_golangdev/internal/transport/rest"
	restHandler "github.com/eskermese/tz_golangdev/internal/transport/rest/handlers"
	"github.com/eskermese/tz_golangdev/pkg/clients/getblock"
	"github.com/eskermese/tz_golangdev/pkg/workers"

	"github.com/eskermese/tz_golangdev/internal/config"
	"github.com/eskermese/tz_golangdev/internal/service"
	"github.com/eskermese/tz_golangdev/pkg/logger"
)

var version = "unknown"

// @title GetBlock TZ
// @version 1.0
// @description Service for getting blocks

// @host localhost:8000
// @BasePath /api/

// Run initializes whole application.
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	log := logger.New("debug", "tz_golangdev")
	defer func() {
		if err := logger.Cleanup(log); err != nil {
			log.Error("error cleanup logs", logger.Error(err))
		}
	}()

	cfg, err := config.New()
	if err != nil {
		log.Fatal("error when initializing the config", logger.Error(err))
	}

	cfg.Version = version

	client := getblock.NewClient(cfg.ApiKey, &http.Client{
		Timeout: 30 * time.Second,
	})

	services := service.New(service.Deps{
		Getblock: client,
	})
	handlers := restHandler.New(restHandler.Deps{
		Service: services,
		Logger:  log,
	})

	g, gCtx := workers.GroupWithContext(ctx)

	g.Go(func() error {
		httpSrv := rest.NewServer(gCtx, cfg, handlers, log)

		if err := httpSrv.Run(); err != nil {
			return fmt.Errorf("http server run: %w", err)
		}

		return nil
	})

	if err = g.Wait(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("error group wait", logger.Error(err))
	}
}
