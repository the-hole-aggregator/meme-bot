package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"meme-bot/internal/adapters/source"
	"meme-bot/internal/adapters/source/downloader"
	"meme-bot/internal/config"
	"meme-bot/internal/util"
	"net/http"
	"time"

	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
	"github.com/mmcdole/gofeed"
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
					_, err := fmt.Scanln(&code)
					if err != nil {
						return "", err
					}
					return code, nil
				}),
			),
			auth.SendCodeOptions{},
		)

		if err := client.Auth().IfNecessary(ctx, flow); err != nil {
			return err
		}

		log.Println("✅ Authorized")

		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		hasher := util.ImagePHasher{}
		downloaderFactory := downloader.DefaultDownloaderFactory{}
		feedParser := gofeed.NewParser()
		// tgSource := source.NewTelegramSource(client, "git_rebase", hasher, downloaderFactory, r)

		// meme, err := tgSource.FetchMeme(ctx)
		// if err != nil {
		// 	return fmt.Errorf("failed to fetch meme: %w", err)
		// }
		rssSource := source.NewRssSource(
			cfg.RSS_SOURCES[2],
			feedParser,
			downloaderFactory,
			r,
			http.DefaultClient,
			hasher,
		)

		meme, err := rssSource.FetchMeme(ctx)
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
