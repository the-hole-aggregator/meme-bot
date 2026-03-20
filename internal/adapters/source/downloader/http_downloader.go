package downloader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
)

type HTTPDownloader struct {
	ctx    context.Context
	url    string
	client *http.Client
}

func NewHTTPDownloader(
	ctx context.Context,
	url string,
	client *http.Client,
) *HTTPDownloader {
	return &HTTPDownloader{
		ctx:    ctx,
		url:    url,
		client: client,
	}
}

func (d *HTTPDownloader) DownloadImage(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer func() {
		closeErr := file.Close()
		if err == nil && closeErr != nil {
			err = fmt.Errorf("close file: %w", closeErr)
		}
	}()

	req, err := http.NewRequestWithContext(d.ctx, http.MethodGet, d.url, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer func() {
		closeErr := resp.Body.Close()
		if err == nil && closeErr != nil {
			err = fmt.Errorf("close body: %w", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("copy body: %w", err)
	}

	return nil
}
