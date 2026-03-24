package ports

import "meme-bot/internal/domain"

type Repository interface {
	GetByID(ID int) (domain.Meme, error)
	Save(meme *domain.Meme) error
	ExistsByHash(hash string) (bool, error)
	GetByStatus(status domain.MemeStatus) ([]domain.Meme, error)
	UpdateStatus(ID int, status domain.MemeStatus) error
	GetOldestApproved() (domain.Meme, error)
	Delete(ID int) error
}
