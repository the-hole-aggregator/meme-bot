package usecase

import (
	"context"
	"fmt"
	"image"
	"log/slog"
	"meme-bot/internal/domain"
	"meme-bot/internal/ports"
	"os"
	"strconv"
)

const maxAttempts = 50

type IngestionUseCase struct {
	repository ports.Repository
	sources    []ports.Source
	logger     *slog.Logger
}

func NewIngestionUseCase(
	repository ports.Repository,
	sources []ports.Source,
	logger *slog.Logger,
) *IngestionUseCase {
	return &IngestionUseCase{
		repository: repository,
		sources:    sources,
		logger:     logger,
	}
}

func (s *IngestionUseCase) Call(ctx context.Context, limit int) error {
	if len(s.sources) == 0 {
		return fmt.Errorf("no sources configured")
	}

	collected := 0
	attempts := 0

	for attempts < maxAttempts && collected < limit {
		attempts++

		source := s.sources[attempts%len(s.sources)]

		meme, filePath, err := source.FetchMeme(ctx)
		if err != nil {
			s.logger.Error(fmt.Errorf("fetch error: %s", err).Error())
			continue
		}

		if !s.validate(meme, filePath) {
			os.Remove(filePath)
			continue
		}

		if err := s.repository.Save(meme); err != nil {
			s.logger.Error(fmt.Errorf("save error: %s", err).Error())
			os.Remove(filePath)
			continue
		}

		collected++
	}

	s.logger.Info("ingestion finished: collected=%d attempts=%d", strconv.Itoa(collected), attempts)

	if collected == 0 {
		return fmt.Errorf("no memes collected after %d attempts", attempts)
	}

	return nil
}

func (s *IngestionUseCase) validate(meme *domain.Meme, filePath string) bool {
	if !isValidImage(filePath) {
		s.logger.Error("image isn't valid")
		return false
	}

	if !isNotEmptyFile(filePath) {
		s.logger.Error("image is empty")
		return false
	}

	exists, err := s.repository.ExistsByHash(meme.PHash)
	if err != nil {
		s.logger.Error(fmt.Errorf("failed on checking existence by phash: %s", err).Error())
		return false
	}

	return !exists
}

func isValidImage(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	_, _, err = image.Decode(f)
	return err == nil
}

func isNotEmptyFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Size() > 0
}
