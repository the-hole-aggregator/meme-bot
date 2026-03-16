package service

import (
	"fmt"
	"log/slog"
	"meme-bot/internal/entity"
)

type PublishService struct {
	publishers []publisher
	repository repository
	logger     slog.Logger
}

func NewPublisherService(publishers []publisher, repository repository) *PublishService {
	return &PublishService{publishers: publishers, repository: repository}
}

func (p *PublishService) Publish() error {
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

		if len(publishErrors) == 3 {
			return err
		}
	}

	return p.repository.UpdateStatus(meme.ID, entity.Posted)
}
