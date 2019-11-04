package config

import (
	"github.com/florinutz/filme/pkg/container"
)

// Config holds flags, env vars and config files contents
type Config struct {
	*container.Container
}

func New() *Config {
	return &Config{
		Container: container.New(),
	}
}
