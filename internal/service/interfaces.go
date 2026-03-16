package service

import "meme-bot/internal/entity"

type repository interface {
	Save(meme entity.Meme) error
	ExistsByHash(hash string) (bool, error)
	GetByStatus(status entity.MemeStatus) ([]entity.Meme, error)
	UpdateStatus(ID string, status entity.MemeStatus) error
	GetOldestApproved() (entity.Meme, error)
}

type source interface {
	FetchMeme() (entity.Meme, error)
}

type publisher interface {
	Publish(meme entity.Meme) error
}

type moderationAPI interface {
	Moderate(meme entity.Meme) error
}
