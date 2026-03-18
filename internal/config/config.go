package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	TG_API_ID   int
	TG_API_HASH string
}

func NewConfig() (cfg *Config, err error) {
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg = &Config{}

	cfg.TG_API_ID, err = getEnvAsInt("TG_API_ID")
	if err != nil {
		return nil, err
	}

	cfg.TG_API_HASH = os.Getenv("TG_API_HASH")

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil

}

func (c *Config) validate() error {
	if c.TG_API_ID == 0 {
		return errors.New("TG_API_ID is required")
	}

	if c.TG_API_HASH == "" {
		return errors.New("TG_API_HASH is required")
	}

	return nil
}

func getEnvAsInt(key string) (int, error) {
	val := os.Getenv(key)
	if val == "" {
		return 0, fmt.Errorf("%s is required", key)
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("%s must be int: %w", key, err)
	}

	return intVal, nil
}
