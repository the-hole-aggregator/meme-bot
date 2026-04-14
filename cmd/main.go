package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"time"

	"meme-bot/internal/adapters"
	"meme-bot/internal/adapters/publisher"
	"meme-bot/internal/adapters/repository"
	"meme-bot/internal/adapters/source"
	"meme-bot/internal/adapters/source/downloader"
	"meme-bot/internal/config"
	"meme-bot/internal/delivery"
	"meme-bot/internal/ports"
	"meme-bot/internal/scheduler"
	usecase "meme-bot/internal/use_case"
	"meme-bot/internal/util"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mmcdole/gofeed"
)

func main() {
	ctx := context.Background()

	// --- CONFIG ---
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	fileRemover := adapters.OSFileRemover{}

	// --- DB ---
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("db hasn't been initialized")
	}
	repo := repository.NewPostgresRepository(pool)

	// --- TELEGRAM BOT ---
	bot, err := tgbotapi.NewBotAPI(cfg.TG_BOT_TOKEN)
	if err != nil {
		log.Fatal(err)
	}

	// --- TELEGRAM CLIENT ---
	if err := os.MkdirAll("tmp", 0755); err != nil {
		log.Fatal(err)
	}

	sessionStorage := &session.FileStorage{
		Path: "tmp/session.json",
	}

	tgClient := telegram.NewClient(cfg.TG_API_ID, cfg.TG_API_HASH, telegram.Options{
		SessionStorage: sessionStorage,
	})

	// --- CHANNELS ---
	ingestionCh := make(chan struct{}, 1)

	// --- PUBLISHERS ---
	moderationPublisher := publisher.NewModerationPublisher(bot, cfg.MODERATION_CHAT_ID)
	tgPublisher := publisher.NewTGPublisher(bot, cfg.TG_CHANNEL_ID)

	// --- USE CASES ---
	ingestionUC := usecase.NewIngestionUseCase(repo, createSources(cfg, tgClient), logger, fileRemover)
	sendToModerationUC := usecase.NewSendToModerationUseCase(moderationPublisher, repo)
	moderationUC := usecase.NewHandleModerationResultUseCase(repo, fileRemover)
	publishUC := usecase.NewPublisherUseCase([]ports.Publisher{tgPublisher}, repo, logger, fileRemover)

	// --- TELEGRAM WORKER ---
	go func() {
		err := tgClient.Run(ctx, func(ctx context.Context) error {
			// AUTH
			flow := auth.NewFlow(
				auth.Constant(
					cfg.PHONE,
					cfg.PASSWORD,
					auth.CodeAuthenticatorFunc(func(ctx context.Context, sentCode *tg.AuthSentCode) (string, error) {
						var code string
						log.Print("Enter code: ")
						_, err := fmt.Scanln(&code)
						return code, err
					}),
				),
				auth.SendCodeOptions{},
			)

			if err := tgClient.Auth().IfNecessary(ctx, flow); err != nil {
				return err
			}

			logger.Info("Telegram client is ready")

			// WORKER LOOP
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()

				case <-ingestionCh:
					logger.Info("Running ingestion...")

					if err := ingestionUC.Call(ctx, 50); err != nil {
						logger.Error("ingestion failed", "err", err)
					}
				}
			}
		})

		if err != nil {
			log.Fatal(err)
		}
	}()

	// --- SCHEDULER ---
	scheduler := scheduler.NewCronScheduler()

	if err := scheduler.RegisterJobs(
		func() {
			logger.Info("Trigger ingestion...")

			select {
			case ingestionCh <- struct{}{}:
			default:
				logger.Info("Ingestion already queued, skipping")
			}
		},
		func() {
			logger.Info("Send memes to moderation...")

			if err := sendToModerationUC.Call(); err != nil {
				logger.Error("send failed", "err", err)
			}
		},
		func() {
			logger.Info("Publishing memes...")

			if err := publishUC.Call(); err != nil {
				logger.Error("publish failed", "err", err)
			}
		},
	); err != nil {
		log.Fatal(err)
	}

	scheduler.Start()
	logger.Info("Scheduler started")

	// --- HANDLERS ---
	moderationHandler := delivery.NewModerationHandler(bot, moderationUC, sendToModerationUC, logger)
	go moderationHandler.Start()

	// --- INITIAL INGESTION ---
	ingestionCh <- struct{}{}

	// --- BLOCK MAIN ---
	select {}
}

func createSources(
	cfg *config.Config,
	client *telegram.Client,
) []ports.Source {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	hasher := util.ImagePHasher{}
	downloaderFactory := downloader.DefaultDownloaderFactory{}
	feedParser := gofeed.NewParser()

	var sources []ports.Source

	for _, src := range cfg.TG_SOURCES {
		sources = append(
			sources,
			source.NewTelegramSource(
				client,
				src,
				hasher,
				downloaderFactory,
				r,
			),
		)
	}

	for _, src := range cfg.RSS_SOURCES {
		sources = append(
			sources,
			source.NewRssSource(
				src,
				feedParser,
				downloaderFactory,
				r,
				http.DefaultClient,
				hasher,
			),
		)
	}

	return sources
}
