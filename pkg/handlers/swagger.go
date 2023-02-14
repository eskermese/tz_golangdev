package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// Swagger для размещения API документации.
type Swagger struct {
	path string
}

func NewSwagger(path string) *Swagger {
	return &Swagger{
		path: path,
	}
}

func (sr Swagger) Path() string {
	return sr.path
}

func (sr Swagger) SetRoutes(r fiber.Router) {
	r.Get(sr.path, swagger.New(swagger.Config{}))
}
