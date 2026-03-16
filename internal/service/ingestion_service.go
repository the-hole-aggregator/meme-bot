package service

import (
	"fmt"
	"meme-bot/internal/entity"
)

// Memes count to fetch.
// 2 for each weekday
const (
	fetchMemesCount = 14
	maxAttempts     = 50
)

type IngestionService struct {
	repository repository
	sources    []source
}

func NewIngestionService(repository repository, sources []source) *IngestionService {
	return &IngestionService{repository: repository, sources: sources}
}

func (s *IngestionService) FetchAndProcess() error {
	collected := 0
	attempts := 0

	for collected < fetchMemesCount && attempts < maxAttempts {
		attempts++

		source := s.sources[attempts%len(s.sources)]

		meme, err := source.FetchMeme()
		if err != nil {
			continue
		}

		if !s.validate(meme) {
			continue
		}

		if err := s.repository.Save(meme); err != nil {
			continue
		}
	}

	if collected < fetchMemesCount {
		return fmt.Errorf(
			"failed to collect enough memes: got %d/%d after %d attempts",
			collected,
			fetchMemesCount,
			attempts,
		)
	}

	return nil
}

func (s *IngestionService) validate(meme entity.Meme) bool {
	exists, err := s.repository.ExistsByHash(meme.Hash)
	if err != nil {
		return false
	}

	return exists
}
