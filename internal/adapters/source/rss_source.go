package source

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"meme-bot/internal/adapters/source/downloader"
	"meme-bot/internal/domain"
	"meme-bot/internal/util"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

type RssSource struct {
	url    string
	parser *gofeed.Parser
	df     downloader.DownloaderFactory
	rand   *rand.Rand
	client *http.Client
	hasher util.Hasher
}

func NewRssSource(
	url string,
	parser *gofeed.Parser,
	df downloader.DownloaderFactory,
	rand *rand.Rand,
	client *http.Client,
	hasher util.Hasher,
) *RssSource {
	return &RssSource{
		url:    url,
		parser: parser,
		df:     df,
		rand:   rand,
		client: client,
		hasher: hasher,
	}
}

func (s *RssSource) FetchMeme(ctx context.Context) (*domain.Meme, error) {
	feed, err := s.parser.ParseURLWithContext(s.url, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to parse rss url %s", err)
	}

	filtered := util.Filter(feed.Items, hasImage)
	if len(filtered) == 0 {
		return nil, errors.New("there are not any items with image")
	}

	item := filtered[s.rand.Intn(len(filtered))]

	filePath := fmt.Sprintf("tmp/%s.jpg", item.GUID)

	imgURL := extractImage(item)
	if imgURL == "" {
		return nil, errors.New("image url can't be empty")
	}

	imageDownloader := s.df.FromURL(ctx, imgURL, s.client)
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
		SourceID:  item.GUID,
		CreatedAt: time.Now(),
	}

	return result, nil
}

var imgRe = regexp.MustCompile(`src="([^"]+)"`)

func hasImage(item *gofeed.Item) bool {
	if media, ok := item.Extensions["media"]; ok {
		if thumbs, ok := media["thumbnail"]; ok && len(thumbs) > 0 {
			return true
		}
	}
	return imgRe.MatchString(item.Description)
}

func extractImage(item *gofeed.Item) string {
	if media, ok := item.Extensions["media"]; ok {
		if thumbs, ok := media["thumbnail"]; ok && len(thumbs) > 0 {
			url := thumbs[0].Attrs["url"]
			return toHD(url)
		}
	}

	return ""
}

func toHD(url string) string {
	url = strings.Replace(url, "preview.redd.it", "i.redd.it", 1)

	if i := strings.Index(url, "?"); i != -1 {
		url = url[:i]
	}

	return url
}
