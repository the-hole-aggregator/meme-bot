package usecase

import (
	"log/slog"
	"meme-bot/internal/domain"
	"meme-bot/internal/ports"
)

type ModerationUseCase struct {
	moderationAPI ports.Moderation
	repository    ports.Repository
	logger        slog.Logger
}

func NewModerationUseCase(moderationAPI ports.Moderation, repository ports.Repository) *ModerationUseCase {
	return &ModerationUseCase{moderationAPI: moderationAPI, repository: repository}
}

func (m *ModerationUseCase) Call() error {
	pendingMemes, err := m.repository.GetByStatus(domain.Pending)
	if err != nil {
		return err
	}

	for _, meme := range pendingMemes {
		err := m.moderationAPI.Moderate(meme)
		if err != nil {
			m.logger.Error(err.Error())
			continue
		}
	}

	return nil
}
