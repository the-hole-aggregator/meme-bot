package usecase

import (
	"errors"
	"meme-bot/internal/domain"
	"meme-bot/internal/ports"
)

type SendToModerationUseCase struct {
	moderation ports.Publisher
	repository ports.Repository
}

func NewSendToModerationUseCase(moderation ports.Publisher, repository ports.Repository) *SendToModerationUseCase {
	return &SendToModerationUseCase{moderation: moderation, repository: repository}
}

func (m *SendToModerationUseCase) Call() error {
	pendingMemes, err := m.repository.GetByStatus(domain.Pending)
	if err != nil {
		return err
	}

	if len(pendingMemes) == 0 {
		return errors.New("there is no any pending meme to moderate")
	}

	if err := m.moderation.Publish(pendingMemes[0]); err != nil {
		return err
	}

	return nil
}
