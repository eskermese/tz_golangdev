package handlers

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/sync/errgroup"
)

type Pinger interface {
	Name() string
	Ping(ctx context.Context) error
}

type Health struct {
	path    string
	pingers []Pinger
}

func NewHealth(path string, pingers ...Pinger) *Health {
	return &Health{path: path, pingers: pingers}
}

func (h Health) SetRoutes(r fiber.Router) {
	r.Get(h.path, h.check)
}

func (h Health) Path() string {
	return h.path
}

func (h Health) check(c *fiber.Ctx) error {
	g, gctx := errgroup.WithContext(c.Context())

	for i := range h.pingers {
		func(i int) {
			g.Go(func() error {
				if err := h.pingers[i].Ping(gctx); err != nil {
					return fmt.Errorf("ping %s : %w", h.pingers[i].Name(), err)
				}

				return nil
			})
		}(i)
	}

	if err := g.Wait(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"detail": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}
