package config

import (
	"github.com/joho/godotenv"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
	"github.com/pkg/errors"
)

type Config struct {
	ARBITER_HOST      string `koanf:"ARBITER_HOST"`
	ROLE              string `koanf:"ROLE"`
	POSTGRES_USER     string `koanf:"POSTGRES_USER"`
	POSTGRES_PASSWORD string `koanf:"POSTGRES_PASSWORD"`
	MASTER_HOST       string `koanf:"MASTER_HOST"`
	SLAVE_HOST        string `koanf:"SLAVE_HOST"`
	PG_DATA           string `koanf:"PG_DATA"`
}

func Load() (config *Config, err error) {
	k := koanf.New(".")

	godotenv.Load()
	err = k.Load(env.Provider("", ".", func(s string) string { return s }), nil)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to load environment variables")
	}
	if err := k.Unmarshal("", &config); err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshall config")
	}

	return config, nil
}
