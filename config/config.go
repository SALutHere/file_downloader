package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `yaml:"env" env-default:"local"`
}

func MustLoad(configPath string) *Config {
	if configPath == "" {
		log.Fatal("config path is not set")
	}

	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("can not read config: %s", err)
	}

	return &cfg
}
