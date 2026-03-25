package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"meme-bot/internal/adapters/publisher"
	"meme-bot/internal/adapters/repository"
	"meme-bot/internal/adapters/source"
	"meme-bot/internal/adapters/source/downloader"
	"meme-bot/internal/config"
	"meme-bot/internal/ports"
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
	"github.com/robfig/cron/v3"
)

func main() {
	ctx := context.Background()

	// --- CONFIG ---
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

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
	sessionStorage := &session.FileStorage{
		Path: "session.json",
	}
	tgClient := telegram.NewClient(cfg.TG_API_ID, cfg.TG_API_HASH, telegram.Options{
		SessionStorage: sessionStorage,
	})

	// --- SOURCES ---
	sources := createSources(cfg, tgClient)

	// --- PUBLISHERS ---
	moderationPublisher := publisher.NewModerationPublisher(bot, cfg.MODERATION_CHAT_ID)
	telegramPublisher := publisher.NewTelegramPublisher(bot, cfg.TG_CHANNEL_ID)

	// --- USECASES ---
	ingestionUC := usecase.NewIngestionUseCase(repo, sources, logger)
	sendToModerationUC := usecase.NewSendToModerationUseCase(moderationPublisher, repo)
	moderationUC := usecase.NewHandleModerationResultUseCase(repo)
	publishUC := usecase.NewPublisherUseCase([]ports.Publisher{telegramPublisher}, repo, logger)

	// -- First start ingestion
	if err := runIngestion(ctx, tgClient, cfg, ingestionUC); err != nil {
		logger.Error("ingestion failed", "err", err)
	}

	// --- CRON ---
	c := cron.New(cron.WithLocation(time.Local))

	// Ingestion: sunday 10:00
	_, err = c.AddFunc("0 10 * * 0", func() {
		logger.Info("Running ingestion...")

		err := runIngestion(ctx, tgClient, cfg, ingestionUC)
		if err != nil {
			logger.Error("ingestion failed", "err", err)
		}
	})
	if err != nil {
		log.Fatal(err)
	}

	// Send to moderation: daily 9:00 and 19:00
	_, err = c.AddFunc("0 9,19 * * *", func() {
		logger.Info("Send memes to moderation...")

		if err := sendToModerationUC.Call(); err != nil {
			logger.Error("send failed", "err", err)
		}
	})
	if err != nil {
		log.Fatal(err)
	}

	// Publish: daily 10:00 and 20:00
	_, err = c.AddFunc("0 10,20 * * *", func() {
		logger.Info("Publishing memes...")

		if err := publishUC.Call(); err != nil {
			logger.Error("publish failed", "err", err)
		}
	})
	if err != nil {
		log.Fatal(err)
	}

	c.Start()

	// --- TELEGRAM CALLBACK HANDLER ---
	go startModerationListener(bot, moderationUC, logger)

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

		return uc.Call(ctx, 14)
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

func startModerationListener(
	bot *tgbotapi.BotAPI,
	uc *usecase.HandleModerationResultUseCase,
	logger *slog.Logger,
) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery == nil {
			continue
		}

		cb := update.CallbackQuery

		parts := strings.Split(cb.Data, ":")
		if len(parts) != 2 {
			continue
		}

		action := parts[0]
		id, _ := strconv.Atoi(parts[1])

		if err := uc.Call(id, usecase.UserSelectionType(action)); err != nil {
			logger.Error("moderation failed", "err", err)
			if _, err := bot.Request(tgbotapi.NewCallback(cb.ID, "ERROR")); err != nil {
				logger.Error(err.Error())
			}
			continue
		}

		if _, err := bot.Request(tgbotapi.NewCallback(cb.ID, "OK")); err != nil {
			logger.Error(err.Error())
		}
	}
}
