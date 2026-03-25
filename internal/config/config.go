package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	PHONE              string
	PASSWORD           string
	TG_BOT_TOKEN       string
	TG_API_ID          int
	TG_API_HASH        string
	TG_SOURCES         []string
	RSS_SOURCES        []string
	DATABASE_URL       string
	MODERATION_CHAT_ID int64
	TG_CHANNEL_ID      int64
}

func NewConfig() (cfg *Config, err error) {
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file ", err)
	}

	cfg = &Config{}

	cfg.PHONE = os.Getenv("PHONE_NUMBER")
	cfg.PASSWORD = os.Getenv("PASSWORD")

	cfg.TG_BOT_TOKEN = os.Getenv("TG_BOT_TOKEN")

	cfg.TG_API_ID, err = getEnvAsInt("TG_API_ID")
	if err != nil {
		return nil, err
	}
	cfg.TG_API_HASH = os.Getenv("TG_API_HASH")

	cfg.TG_SOURCES = strings.Split(os.Getenv("TG_SOURCES"), ",")
	for i := range cfg.TG_SOURCES {
		cfg.TG_SOURCES[i] = strings.TrimSpace(cfg.TG_SOURCES[i])
	}
	cfg.RSS_SOURCES = strings.Split(os.Getenv("RSS_SOURCES"), ",")
	for i := range cfg.RSS_SOURCES {
		cfg.RSS_SOURCES[i] = strings.TrimSpace(cfg.RSS_SOURCES[i])
	}

	cfg.DATABASE_URL = os.Getenv(("DATABASE_URL"))

	cfg.MODERATION_CHAT_ID, err = getEnvAsInt64("MODERATION_CHAT_ID")
	if err != nil {
		return nil, err
	}

	cfg.TG_CHANNEL_ID, err = getEnvAsInt64("TG_CHANNEL_ID")
	if err != nil {
		return nil, err
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil

}

func (c *Config) validate() error {
	if c.TG_BOT_TOKEN == "" {
		return errors.New("TG_BOT_TOKEN is required")
	}

	if c.TG_API_ID == 0 {
		return errors.New("TG_API_ID is required")
	}

	if c.TG_API_HASH == "" {
		return errors.New("TG_API_HASH is required")
	}

	if c.PHONE == "" {
		return errors.New("PHONE_NUMBER is required")
	}

	if c.PASSWORD == "" {
		return errors.New("PASSWORD is required")
	}

	if len(c.TG_SOURCES) == 0 {
		return errors.New("TG_SOURCES can't be empty")
	}

	if len(c.RSS_SOURCES) == 0 {
		return errors.New("RSS_SOURCES can't be empty")
	}

	if c.DATABASE_URL == "" {
		return errors.New("DATABASE_URL is required")
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

func getEnvAsInt64(key string) (int64, error) {
	val := os.Getenv(key)
	if val == "" {
		return 0, fmt.Errorf("%s is required", key)
	}

	intVal, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s must be int64: %w", key, err)
	}

	return intVal, nil
}
