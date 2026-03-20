package ports

import "meme-bot/internal/domain"

type Repository interface {
	Save(meme *domain.Meme) error
	ExistsByHash(hash string) (bool, error)
	GetByStatus(status domain.MemeStatus) ([]domain.Meme, error)
	UpdateStatus(ID int, status domain.MemeStatus) error
	GetOldestApproved() (domain.Meme, error)
}
