package ports

import "meme-bot/internal/domain"

type Moderation interface {
	Moderate(meme domain.Meme) error
}
