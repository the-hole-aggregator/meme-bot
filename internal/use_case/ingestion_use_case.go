package service

import (
	"fmt"
	"log/slog"
	"meme-bot/internal/domain"
	"meme-bot/internal/ports"
)

const maxAttempts = 50

type IngestionUseCase struct {
	repository ports.Repository
	sources    []ports.Source
	logger     slog.Logger
}

func NewIngestionUseCase(repository ports.Repository, sources []ports.Source) *IngestionUseCase {
	return &IngestionUseCase{
		repository: repository,
		sources:    sources,
	}
}

func (s *IngestionUseCase) Call(limit int) error {
	if len(s.sources) == 0 {
		return fmt.Errorf("no sources configured")
	}

	collected := 0
	attempts := 0

	for attempts < maxAttempts && collected < limit {
		attempts++

		source := s.sources[attempts%len(s.sources)]

		meme, err := source.FetchMeme()
		if err != nil {
			s.logger.Error(fmt.Errorf("fetch error: %s", err).Error())
			continue
		}

		if !s.validate(meme) {
			continue
		}

		if err := s.repository.Save(meme); err != nil {
			s.logger.Error(fmt.Errorf("save error: %s", err).Error())
			continue
		}

		collected++
	}

	s.logger.Error(fmt.Errorf("ingestion finished: collected=%d attempts=%d", collected, attempts).Error())

	if collected == 0 {
		return fmt.Errorf("no memes collected after %d attempts", attempts)
	}

	return nil
}

func (s *IngestionUseCase) validate(meme domain.Meme) bool {
	exists, err := s.repository.ExistsByHash(meme.Hash)
	if err != nil {
		return false
	}

	return !exists
}
