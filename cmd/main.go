package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"meme-bot/internal/adapters/repository"
	"meme-bot/internal/adapters/source"
	"meme-bot/internal/adapters/source/downloader"
	"meme-bot/internal/config"
	"meme-bot/internal/ports"
	usecase "meme-bot/internal/use_case"
	"meme-bot/internal/util"
	"net/http"
	"os"
	"time"

	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
	"github.com/jackc/pgx/v5/pgxpool"
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

		logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

		pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
		if err != nil {
			log.Fatal("db hasn't been initialized")
		}
		repository := repository.NewPostgresRepository(pool)
		sources := createSources(cfg, client)

		ingestionUseCase := usecase.NewIngestionUseCase(
			repository,
			sources,
			logger,
		)

		if err := ingestionUseCase.Call(ctx, 14); err != nil {
			log.Fatalf("failed on fetching memes %s", err)
		}

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

func createSources(cfg *config.Config, client *telegram.Client) []ports.Source {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	hasher := util.ImagePHasher{}
	downloaderFactory := downloader.DefaultDownloaderFactory{}
	feedParser := gofeed.NewParser()

	sources := make([]ports.Source, 0)

	for _, src := range cfg.TG_SOURCES {
		tgSource := source.NewTelegramSource(client, src, hasher, downloaderFactory, r)
		sources = append(sources, tgSource)
	}

	for _, src := range cfg.RSS_SOURCES {
		rssSource := source.NewRssSource(
			src,
			feedParser,
			downloaderFactory,
			r,
			http.DefaultClient,
			hasher,
		)
		sources = append(sources, rssSource)
	}

	return sources
}
