package downloader

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gotd/td/telegram/downloader"
	"github.com/gotd/td/tg"
)

type TgImageDownloader struct {
	ctx context.Context
	api *tg.Client
	msg *tg.Message
}

func NewTgImageDownloader(ctx context.Context, api *tg.Client, msg *tg.Message) *TgImageDownloader {
	return &TgImageDownloader{
		ctx: ctx,
		api: api,
		msg: msg,
	}
}

func (tgd *TgImageDownloader) DownloadImage(filePath string) error {
	media, ok := tgd.msg.Media.(*tg.MessageMediaPhoto)
	if !ok {
		return fmt.Errorf("not a photo")
	}

	photo, ok := media.Photo.(*tg.Photo)
	if !ok {
		return fmt.Errorf("invalid photo")
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

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Print("error on closing file ", err)
		}

	}()

	d := downloader.NewDownloader()

	location := &tg.InputPhotoFileLocation{
		ID:            photo.ID,
		AccessHash:    photo.AccessHash,
		FileReference: photo.FileReference,
		ThumbSize:     bestType,
	}

	builder := d.Download(tgd.api, location)

	_, err = builder.ToPath(tgd.ctx, filePath)
	if err != nil {
		return err
	}

	return nil
}
