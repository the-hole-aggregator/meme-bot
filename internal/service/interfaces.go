package service

import "meme-bot/internal/entity"

type repository interface {
	Save(meme entity.Meme) error
	ExistsByHash(hash string) (bool, error)
	GetByStatus(status entity.MemeStatus) ([]entity.MemeStatus, error)
}

type source interface {
	FetchMeme() (entity.Meme, error)
}

type publisher interface {
	Publish(meme entity.Meme) error
}
