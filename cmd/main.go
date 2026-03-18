package main

import (
	"context"
	"fmt"
	"log"
	"meme-bot/internal/adapters/source"
	"meme-bot/internal/config"
	"time"

	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
)

func main() {
	ctx := context.Background()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	sessionStorage := &session.FileStorage{
		Path: "session.json",
	}

	client := telegram.NewClient(cfg.TG_API_ID, cfg.TG_API_HASH, telegram.Options{
		SessionStorage: sessionStorage,
	})

	err = client.Run(ctx, func(ctx context.Context) error {

		flow := auth.NewFlow(
			auth.Constant(
				cfg.PHONE,
				cfg.PASSWORD,
				auth.CodeAuthenticatorFunc(func(ctx context.Context, sentCode *tg.AuthSentCode) (string, error) {
					var code string
					fmt.Print("Enter code: ")
					fmt.Scanln(&code)
					return code, nil
				}),
			),
			auth.SendCodeOptions{},
		)

		if err := client.Auth().IfNecessary(ctx, flow); err != nil {
			return err
		}

		log.Println("✅ Authorized")

		tgSource := source.NewTelegramSource(client, "tupi4ek_degradanta")

		meme, err := tgSource.FetchMeme(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch meme: %w", err)
		}

		log.Printf("🎉 Meme fetched: %+v\n", meme)

		// TODO(egrischenkov): change to scheduler
		for {
			time.Sleep(10 * time.Minute)
			log.Println("alive...")
		}
	})

	if err != nil {
		log.Fatal(err)
	}
}
