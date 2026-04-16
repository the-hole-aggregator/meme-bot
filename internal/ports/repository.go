package ports

import (
	"errors"
	"meme-bot/internal/domain"
)

var ErrMemesEnded = errors.New("memes have ended")

type Repository interface {
	GetByID(ID int) (domain.Meme, error)
	Save(meme *domain.Meme) error
	ExistsByHash(hash string) (bool, error)
	ExistsBySourceID(sourceID string) (bool, error)
	GetByStatus(status domain.MemeStatus) ([]domain.Meme, error)
	UpdateStatus(ID int, status domain.MemeStatus) error
	GetOldestApproved() (domain.Meme, error)
	GetOldestPending() (domain.Meme, error)
	Delete(ID int) error
}
