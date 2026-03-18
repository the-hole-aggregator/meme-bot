package ports

import "meme-bot/internal/domain"

type Source interface {
	FetchMeme() (domain.Meme, error)
}
