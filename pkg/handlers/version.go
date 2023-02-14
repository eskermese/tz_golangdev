package handlers

import (
	"github.com/gofiber/fiber/v2"
)

const APIVersion = "v1"

// VersionResponse - ответ на запрос версии.
type VersionResponse struct {
	API     string `json:"api"`
	Version string `json:"version"`
}

type Version struct {
	version string
	path    string
}

func NewVersion(path string, version string) *Version {
	return &Version{version: version, path: path}
}

func (v Version) SetRoutes(r fiber.Router) {
	r.Get(v.path, versionHandlerFunc(v.version))
}

func versionHandlerFunc(version string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		vr := VersionResponse{
			API:     APIVersion,
			Version: version,
		}

		return c.Status(fiber.StatusOK).JSON(vr)
	}
}
