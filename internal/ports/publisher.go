package ports

import "meme-bot/internal/domain"

type Publisher interface {
	Publish(meme domain.Meme) error
}
