package usecase

import (
	"meme-bot/internal/ports"
)

type SendToModerationUseCase struct {
	moderation ports.Publisher
	repo       ports.Repository
}

func NewSendToModerationUseCase(moderation ports.Publisher, repository ports.Repository) *SendToModerationUseCase {
	return &SendToModerationUseCase{moderation: moderation, repo: repository}
}

func (uc *SendToModerationUseCase) Call() error {
	meme, err := uc.repo.GetOldestPending()
	if err != nil {
		return err
	}

	if err := uc.moderation.Publish(meme); err != nil {
		return err
	}

	return nil
}
