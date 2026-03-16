package service

import (
	"log/slog"
	"meme-bot/internal/entity"
)

type ModerationService struct {
	moderationAPI moderationAPI
	repository    repository
	logger        slog.Logger
}

func NewModerationService(moderationAPI moderationAPI, repository repository) *ModerationService {
	return &ModerationService{moderationAPI: moderationAPI, repository: repository}
}

func (m *ModerationService) Moderate() error {
	pendingMemes, err := m.repository.GetByStatus(entity.Pending)
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
