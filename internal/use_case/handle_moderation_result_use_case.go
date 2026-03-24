package usecase

import (
	"fmt"
	"meme-bot/internal/domain"
	"meme-bot/internal/ports"
	"os"

	"github.com/go-faster/errors"
)

type UserSelectionType string

const (
	Approved UserSelectionType = "approved"
	Rejected UserSelectionType = "rejected"
)

type HandleModerationResultUseCase struct {
	repo ports.Repository
}

func NewHandleModerationResultUseCase(repo ports.Repository) *HandleModerationResultUseCase {
	return &HandleModerationResultUseCase{repo: repo}
}

func (h *HandleModerationResultUseCase) Call(id int, userSelection UserSelectionType) error {
	switch userSelection {
	case Approved:
		if err := h.repo.UpdateStatus(id, domain.Approved); err != nil {
			return err
		}
	case Rejected:
		meme, err := h.repo.GetByID(id)
		if err != nil {
			return errors.Wrap(err, "failed on getting meme by ID")
		}

		if err := os.Remove(fmt.Sprintf("tmp/%s.jpg", meme.SourceID)); err != nil {
			return err
		}

		if err := h.repo.Delete(id); err != nil {
			return err
		}
	default:
		return fmt.Errorf("there is no such user selection type: %s", userSelection)
	}

	return nil
}
