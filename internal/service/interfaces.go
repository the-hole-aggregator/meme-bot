package service

import "meme-bot/internal/entity"

type repository interface {
	save(meme entity.Meme) error
	existsByHash(hash string) (bool, error)
	getByStatus(status entity.MemeStatus) ([]entity.MemeStatus, error)
}

type source interface {
	fetchMemes() ([]entity.Meme, error)
}

type publisher interface {
	publish(meme entity.Meme) error
}
