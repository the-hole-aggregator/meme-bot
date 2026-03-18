package main

import (
	"context"
	"log"
	"meme-bot/internal/adapters/source"
	"meme-bot/internal/config"
	"time"

	"github.com/gotd/td/telegram"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	tgClient := telegram.NewClient(cfg.TG_API_ID, cfg.TG_API_HASH, telegram.Options{})
	if err := tgClient.Run(ctx, func(ctx context.Context) error {
		log.Println("Telegram client initialized")
		return nil
	}); err != nil {
		log.Fatalf("failed to init telegram client: %v", err)
	}

	tgSource := source.NewTelegramSource(tgClient, "tupi4ek_degradanta")
	tgSource.FetchMeme(ctx)
}
