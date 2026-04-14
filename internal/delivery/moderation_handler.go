package delivery

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	usecase "meme-bot/internal/use_case"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
)

type ModerationHandler struct {
	bot                *tgbotapi.BotAPI
	handleModerationUC *usecase.HandleModerationResultUseCase
	sendToModerationUC *usecase.SendToModerationUseCase
	logger             *slog.Logger
}

func NewModerationHandler(
	bot *tgbotapi.BotAPI,
	handleModerationUC *usecase.HandleModerationResultUseCase,
	sendToModerationUC *usecase.SendToModerationUseCase,
	logger *slog.Logger,
) *ModerationHandler {
	return &ModerationHandler{
		bot:                bot,
		handleModerationUC: handleModerationUC,
		sendToModerationUC: sendToModerationUC,
		logger:             logger,
	}
}

func (h *ModerationHandler) Start(runIngestion func()) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := h.bot.GetUpdatesChan(u)

	for update := range updates {
		if err := h.handleUpdate(update); err != nil {
			h.logger.Error(err.Error())

			if errors.Is(err, pgx.ErrNoRows) {
				runIngestion()
			}
		}
	}
}

func (h *ModerationHandler) handleUpdate(update tgbotapi.Update) error {
	if update.CallbackQuery == nil {
		return errors.New("callbackQuery is nil")
	}

	cb := update.CallbackQuery

	parts := strings.Split(cb.Data, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid moderation data format: %v", parts)
	}

	action := parts[0]
	id, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid moderation id format %v", err)
	}

	h.logger.Info("handle moderation event with", "action", action, "id", id)

	if err := h.handleModerationUC.Call(id, usecase.UserSelectionType(action)); err != nil {
		h.logger.Error("moderation failed", "err", err)
		h.answerCallback(cb.ID, "ERROR")

		return err
	}

	if usecase.UserSelectionType(action) == usecase.Rejected {
		if err := h.sendToModerationUC.Call(); err != nil {
			h.logger.Error("failed to send additional meme for moderation", "err", err)
			h.answerCallback(cb.ID, "ERROR")

			return err
		}
	}

	h.answerCallback(cb.ID, "OK")

	return nil
}

func (h *ModerationHandler) answerCallback(id, text string) {
	if _, err := h.bot.Request(tgbotapi.NewCallback(id, text)); err != nil {
		h.logger.Error("callback failed", "err", err)
	}
}
