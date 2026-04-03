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

	"github.com/go-faster/errors"
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

	// --- TELEGRAM BOT (moderation) ---
	bot, err := tgbotapi.NewBotAPI(cfg.TG_BOT_TOKEN)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to create bot api"))
	}

	// --- TELEGRAM CLIENT (meme's parsing) ---
	if err := os.MkdirAll("tmp", 0755); err != nil {
		log.Fatal(errors.Wrap(err, "failed to create tmp directory"))	
	}
	sessionStorage := &session.FileStorage{
		Path: "tmp/session.json",
		
	}
	tgClient := telegram.NewClient(cfg.TG_API_ID, cfg.TG_API_HASH, telegram.Options{
		SessionStorage: sessionStorage,
	})
	
	// --- PUBLISHERS ---
	moderationPublisher := publisher.NewModerationPublisher(bot, cfg.MODERATION_CHAT_ID)
	tgPublisher := publisher.NewTGPublisher(bot, cfg.TG_CHANNEL_ID)

	// --- USE CASES ---
	ingestionUC := usecase.NewIngestionUseCase(repo, createSources(cfg, tgClient), logger, fileRemover)
	sendToModerationUC := usecase.NewSendToModerationUseCase(moderationPublisher, repo)
	moderationUC := usecase.NewHandleModerationResultUseCase(repo, fileRemover)
	publishUC := usecase.NewPublisherUseCase([]ports.Publisher{tgPublisher}, repo, logger, fileRemover)

	// Init ingestion
	logger.Info("Running ingestion...")
	if err := runIngestion(ctx, tgClient, cfg, ingestionUC); err != nil {
		logger.Error("ingestion failed", "err", err)
	}

	scheduler := scheduler.NewCronScheduler()
	logger.Info("Scheduler has been initialized...")
	logger.Info("Current time: %v, UTC: %v\n", time.Now().String(), time.Now().Local)

	if err := scheduler.RegisterJobs(
		func() {
			logger.Info("Running ingestion...")

			err := runIngestion(ctx, tgClient, cfg, ingestionUC)
			if err != nil {
				logger.Error("ingestion failed", "err", err)
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
	logger.Info("Scheduler has been started...")
	
	// --- HANDLERS ---
	moderationHandler := delivery.NewModerationHandler(bot, moderationUC, sendToModerationUC, logger)

	// --- TELEGRAM CALLBACK HANDLER ---
	go moderationHandler.Start()

	// --- BLOCK MAIN ---
	select {}
}

func runIngestion(
	ctx context.Context,
	client *telegram.Client,
	cfg *config.Config,
	uc *usecase.IngestionUseCase,
) error {

	return client.Run(ctx, func(ctx context.Context) error {

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

		if err := client.Auth().IfNecessary(ctx, flow); err != nil {
			return err
		}

		return uc.Call(ctx, 20)
	})
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
		sources = append(sources,
			source.NewTelegramSource(client, src, hasher, downloaderFactory, r))
	}

	for _, src := range cfg.RSS_SOURCES {
		sources = append(sources,
			source.NewRssSource(
				src,
				feedParser,
				downloaderFactory,
				r,
				http.DefaultClient,
				hasher,
			))
	}

	return sources
}
