package downloader

import (
	"context"
	"net/http"

	"github.com/gotd/td/tg"
)

type DownloaderFactory interface {
	FromTelegram(ctx context.Context, api *tg.Client, msg *tg.Message) ImageDownloader
	FromURL(ctx context.Context, url string, client *http.Client) ImageDownloader
}

type ImageDownloader interface {
	DownloadImage(filePath string) error
}

type DefaultDownloaderFactory struct{}

func (f DefaultDownloaderFactory) FromTelegram(ctx context.Context, api *tg.Client, msg *tg.Message) ImageDownloader {
	return NewTgImageDownloader(ctx, api, msg)
}

func (f DefaultDownloaderFactory) FromURL(ctx context.Context, url string, client *http.Client) ImageDownloader {
	return NewHTTPDownloader(ctx, url, client)
}
