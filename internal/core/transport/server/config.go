package core_http_server

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Addr            string        `envconfig:"ADDR" required:"true"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" required:"true"`
}


func NewHTTPConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("HTTP", &config); err != nil{
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}

func NewHTTPConfigMust() Config {
	cfg, err := NewHTTPConfig()

	if err != nil{
		err = fmt.Errorf("get HTTP server config: %w", err)
		panic(err)
	}
	return cfg
}