package restHandler

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	loggerMW "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/eskermese/tz_golangdev/docs"
	"github.com/eskermese/tz_golangdev/internal/config"
	"github.com/eskermese/tz_golangdev/internal/core"
	"github.com/eskermese/tz_golangdev/pkg/handlers"
	"github.com/eskermese/tz_golangdev/pkg/logger"
)

type Service interface {
	GetMaxTransactionChange(ctx context.Context) (core.Transaction, error)
}

type Deps struct {
	Service Service
	Logger  logger.Logger
}

type Handler struct {
	service Service
	logger  logger.Logger
}

func New(deps Deps) *Handler {
	return &Handler{
		logger:  deps.Logger,
		service: deps.Service,
	}
}

func (h *Handler) InitRouter(cfg *config.Config, router fiber.Router) {
	router.Use(
		recover.New(),
		loggerMW.New(),
	)

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", cfg.ClientHTTP.Host, cfg.ClientHTTP.Port)
	if cfg.Environment != config.EnvLocal {
		docs.SwaggerInfo.Host = cfg.ClientHTTP.Host
	}

	handlers.NewHealth("/health").SetRoutes(router)
	handlers.NewSwagger("/swagger/*").SetRoutes(router)
	handlers.NewVersion("/version", cfg.Version).SetRoutes(router)

	h.initAPI(router)
}

func (h *Handler) initAPI(router fiber.Router) {
	api := router.Group("/api/v1")
	{
		h.initTransactionRoutes(api)
	}
}

type response struct {
	Detail string `json:"detail"`
}
