package ports

import (
	"context"
	"meme-bot/internal/domain"
)

type Source interface {
	FetchMeme(ctx context.Context) (*domain.Meme, error)
}
