package core_logger

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)


type Config struct {
	LogLevel string `envconfig:"LEVEL" required:"true"`
	Path     string	`envconfig:"PATH" required:"true"`
}


func NewLoggerConfig() (Config, error){
	var config Config

	if err := envconfig.Process("LOGGER", &config); err != nil{
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config{
	config, err := NewLoggerConfig()
	if err != nil{
		err = fmt.Errorf("get logger config: %w", err)
		panic(err)
	}
	return config
}