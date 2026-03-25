package usecase

import (
	"fmt"
	"log/slog"
	"meme-bot/internal/domain"
	"meme-bot/internal/ports"
	"os"
)

type PublishUseCase struct {
	publishers []ports.Publisher
	repo       ports.Repository
	logger     *slog.Logger
}

func NewPublisherUseCase(publishers []ports.Publisher, repository ports.Repository, logger *slog.Logger) *PublishUseCase {
	return &PublishUseCase{publishers: publishers, repo: repository, logger: logger}
}

func (uc *PublishUseCase) Call() error {
	meme, err := uc.repo.GetOldestApproved()
	if err != nil {
		return err
	}

	var publishErrors []error

	for _, publisher := range uc.publishers {
		if err := publisher.Publish(meme); err != nil {
			publishErrors = append(publishErrors, err)
		}
	}

	if len(publishErrors) > 0 {
		err := fmt.Errorf("publish errors: %v", publishErrors)
		uc.logger.Error(err.Error())

		if len(publishErrors) == len(uc.publishers) {
			return err
		}
	}

	if err := uc.repo.Delete(meme.ID); err != nil {
		return err
	}

	if err := os.Remove(fmt.Sprintf("tmp/%s.jpg", meme.SourceID)); err != nil {
		return err
	}

	return uc.repo.UpdateStatus(meme.ID, domain.Posted)
}
