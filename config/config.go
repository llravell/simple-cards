package config

import (
	"errors"
	"flag"

	"github.com/caarlos0/env"
)

const (
	_defaultAddr        = ":8080"
	_defaultDatabaseURI = ""
	_defaultJWTSecret   = "secret"
)

var ErrEmptyDatabaseURI = errors.New("got empty database uri")

type Config struct {
	Addr        string `env:"RUN_ADDRESS"`
	DatabaseURI string `env:"DATABASE_URI"`
	JWTSecret   string `env:"JWT_SECRET"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{
		Addr:        _defaultAddr,
		DatabaseURI: _defaultDatabaseURI,
		JWTSecret:   _defaultJWTSecret,
	}

	err := env.Parse(cfg)
	if err != nil {
		return cfg, err
	}

	flag.StringVar(&cfg.Addr, "a", cfg.Addr, "Server address as host:port")
	flag.StringVar(&cfg.DatabaseURI, "d", cfg.DatabaseURI, "Base address for redirect as host:port")
	flag.Parse()

	return cfg, cfg.Validate()
}

func (c *Config) Validate() error {
	if c.DatabaseURI == "" {
		return ErrEmptyDatabaseURI
	}

	return nil
}
