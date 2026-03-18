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

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/downloader"
	"github.com/gotd/td/tg"
)

type TelegramSource struct {
	client  *telegram.Client
	channel string
	limit   int
}

func NewTelegramSource(client *telegram.Client, channel string) *TelegramSource {
	return &TelegramSource{
		client:  client,
		channel: channel,
		limit:   50,
	}
}

func (s *TelegramSource) FetchMeme(ctx context.Context) (*domain.Meme, error) {
	var result *domain.Meme

	err := s.client.Run(ctx, func(ctx context.Context) error {
		api := s.client.API()

		resolved, err := api.ContactsResolveUsername(ctx, &tg.ContactsResolveUsernameRequest{
			Username: s.channel,
		})
		if err != nil {
			return fmt.Errorf("failed to resolve user name %s", err)
		}

		ch, ok := resolved.Chats[0].(*tg.Channel)
		if !ok {
			return errors.New("not a channel")
		}

		history, err := api.MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{
			Peer: &tg.InputPeerChannel{
				ChannelID:  ch.ID,
				AccessHash: ch.AccessHash,
			},
			Limit: s.limit,
		})
		if err != nil {
			return err
		}

		historyMessages, ok := history.AsModified()
		if !ok {
			return fmt.Errorf("can't map history messages")
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
			return errors.New("no media found")
		}

		rand.NewSource(time.Now().UnixNano())
		m := mediaMessages[rand.Intn(len(mediaMessages))]

		filePath, err := s.downloadPhoto(ctx, api, m)
		if err != nil {
			return err
		}

		hash, err := util.ComputePHash(filePath)
		if err != nil {
			os.Remove(filePath)
			return err
		}

		result = &domain.Meme{
			PHash:     hash.ToString(),
			Status:    domain.Pending,
			Source:    domain.Telegram,
			SourceID:  fmt.Sprintf("%d", m.ID),
			CreatedAt: time.Now(),
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *TelegramSource) downloadPhoto(
	ctx context.Context,
	api *tg.Client,
	msg *tg.Message,
) (string, error) {

	media, ok := msg.Media.(*tg.MessageMediaPhoto)
	if !ok {
		return "", fmt.Errorf("not a photo")
	}

	photo, ok := media.Photo.(*tg.Photo)
	if !ok {
		return "", fmt.Errorf("invalid photo")
	}

	var bestType string
	maxSize := 0

	for _, size := range photo.Sizes {
		switch s := size.(type) {

		case *tg.PhotoSize:
			if s.Size > maxSize {
				maxSize = s.Size
				bestType = s.Type
			}

		case *tg.PhotoSizeProgressive:
			if len(s.Sizes) > 0 {
				last := s.Sizes[len(s.Sizes)-1]
				if last > maxSize {
					maxSize = last
					bestType = s.Type
				}
			}
		}
	}

	filePath := fmt.Sprintf("tmp/%d.jpg", msg.ID)

	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	d := downloader.NewDownloader()

	location := &tg.InputPhotoFileLocation{
		ID:            photo.ID,
		AccessHash:    photo.AccessHash,
		FileReference: photo.FileReference,
		ThumbSize:     bestType,
	}

	builder := d.Download(api, location)

	_, err = builder.ToPath(ctx, filePath)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
