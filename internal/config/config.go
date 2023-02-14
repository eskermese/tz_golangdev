package config

import (
	"fmt"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
)

const (
	EnvLocal = "local"
	Prod     = "prod"
)

type Config struct {
	Version string
	*Servers
	*GetBlockConfigs
}

type Servers struct {
	ClientHTTP *ClientServiceConfigs

	Environment string `long:"env" env:"ENVIRONMENT" description:"app environment" default:"local"`
	JWTKey      string `long:"jwt-key" env:"JWT_KEY" description:"JWT secret key" required:"false" default:"some-secret"`
}

type ClientServiceConfigs struct {
	Schema             string        `long:"schema" env:"HTTP_SCHEMA" description:"HTTP url schema" required:"false" default:"http"`
	Host               string        `long:"host" env:"HTTP_HOST" description:"HTTP host" required:"false" default:"localhost"`
	Port               int           `long:"port" env:"HTTP_PORT" description:"HTTP port" required:"false" default:"8000"`
	ReadTimeout        time.Duration `long:"read-timeout" env:"HTTP_READ_TIMEOUT" description:"HTTP read timeout" required:"false" default:"20s"`
	WriteTimeout       time.Duration `long:"write-timeout" env:"HTTP_WRITE_TIMEOUT" description:"HTTP write timeout" required:"false" default:"20s"`
	MaxHeaderMegabytes int           `long:"max-header-mg" env:"HTTP_MAX_HEADER_MEGABYTES" description:"HTTP max header mg" required:"false" default:"1"`
}

type GetBlockConfigs struct {
	ApiKey string `long:"getblock-apikey" env:"GETBLOCK_APIKEY" description:"Api key for access to GetBlock" required:"true"`
}

func New() (*Config, error) {
	defer os.Clearenv()

	c := &Config{}
	p := flags.NewParser(c, flags.Default|flags.IgnoreUnknown)

	if _, err := p.Parse(); err != nil {
		return nil, fmt.Errorf("error parsing config options: %w", err)
	}

	return c, nil
}
