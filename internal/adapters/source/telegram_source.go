package source

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"meme-bot/internal/domain"
	"meme-bot/internal/util"

	tgdownloader "meme-bot/internal/adapters/source/downloader"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
)

type TelegramSource struct {
	client  *telegram.Client
	channel string
	limit   int
	hasher  util.Hasher
	df      tgdownloader.DownloaderFactory
	rand    *rand.Rand
}

func NewTelegramSource(
	client *telegram.Client,
	channel string,
	hasher util.Hasher,
	df tgdownloader.DownloaderFactory,
	rand *rand.Rand,
) *TelegramSource {
	return &TelegramSource{
		client:  client,
		channel: channel,
		limit:   100,
		hasher:  hasher,
		df:      df,
		rand:    rand,
	}
}

func (s *TelegramSource) FetchMeme(ctx context.Context) (*domain.Meme, error) {
	api := s.client.API()

	resolved, err := api.ContactsResolveUsername(ctx, &tg.ContactsResolveUsernameRequest{
		Username: s.channel,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to resolve user name %s", err)
	}

	tgChan, ok := resolved.Chats[0].(*tg.Channel)
	if !ok {
		return nil, errors.New("not a channel")
	}

	history, err := api.MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{
		Peer: &tg.InputPeerChannel{
			ChannelID:  tgChan.ID,
			AccessHash: tgChan.AccessHash,
		},
		Limit: s.limit,
	})
	if err != nil {
		return nil, err
	}

	historyMessages, ok := history.AsModified()
	if !ok {
		return nil, fmt.Errorf("can't map history messages")
	}

	var mediaMessages []*tg.Message

	for _, msg := range historyMessages.GetMessages() {
		m, ok := msg.(*tg.Message)
		if !ok {
			continue
		}

		if m.Media != nil {
			mediaMessages = append(mediaMessages, m)
		}
	}

	if len(mediaMessages) == 0 {
		return nil, errors.New("no media found")
	}

	m := mediaMessages[rand.Intn(len(mediaMessages))]
	filePath := fmt.Sprintf("tmp/%d.jpg", m.ID)

	imageDownloader := s.df.FromTelegram(ctx, api, m)
	if err := imageDownloader.DownloadImage(filePath); err != nil {
		return nil, err
	}

	hash, err := s.hasher.ComputePHash(filePath)
	if err != nil {

		fileErr := os.Remove(filePath)
		if fileErr != nil {
			return nil, fmt.Errorf("failed on computing hash: %s %s", err, fileErr)
		}

		return nil, err
	}

	result := &domain.Meme{
		PHash:     hash,
		Status:    domain.Pending,
		Source:    domain.Telegram,
		SourceID:  fmt.Sprintf("%d", m.ID),
		CreatedAt: time.Now(),
	}

	return result, nil
}
