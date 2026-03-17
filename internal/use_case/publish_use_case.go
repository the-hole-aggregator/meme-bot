package service

import (
	"fmt"
	"log/slog"
	"meme-bot/internal/domain"
	"meme-bot/internal/ports"
)

type PublishUseCase struct {
	publishers []ports.Publisher
	repository ports.Repository
	logger     slog.Logger
}

func NewPublisherUseCase(publishers []ports.Publisher, repository ports.Repository) *PublishUseCase {
	return &PublishUseCase{publishers: publishers, repository: repository}
}

func (p *PublishUseCase) Call() error {
	meme, err := p.repository.GetOldestApproved()
	if err != nil {
		return err
	}

	var publishErrors []error

	for _, publisher := range p.publishers {
		if err := publisher.Publish(meme); err != nil {
			publishErrors = append(publishErrors, err)
		}
	}

	if len(publishErrors) > 0 {
		err := fmt.Errorf("publish errors: %v", publishErrors)
		p.logger.Error(err.Error())

		if len(publishErrors) == len(p.publishers) {
			return err
		}
	}

	return p.repository.UpdateStatus(meme.ID, domain.Posted)
}
