package delivery

import (
	"log/slog"
	"strconv"
	"strings"

	usecase "meme-bot/internal/use_case"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ModerationHandler struct {
	bot    *tgbotapi.BotAPI
	uc     *usecase.HandleModerationResultUseCase
	logger *slog.Logger
}

func NewModerationHandler(
	bot *tgbotapi.BotAPI,
	uc *usecase.HandleModerationResultUseCase,
	logger *slog.Logger,
) *ModerationHandler {
	return &ModerationHandler{
		bot:    bot,
		uc:     uc,
		logger: logger,
	}
}

func (h *ModerationHandler) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := h.bot.GetUpdatesChan(u)

	for update := range updates {
		h.handleUpdate(update)
	}
}

func (h *ModerationHandler) handleUpdate(update tgbotapi.Update) {
	if update.CallbackQuery == nil {
		return
	}

	cb := update.CallbackQuery

	parts := strings.Split(cb.Data, ":")
	if len(parts) != 2 {
		return
	}

	action := parts[0]
	id, err := strconv.Atoi(parts[1])
	if err != nil {
		h.logger.Error("invalid id", "err", err)
		return
	}

	if err := h.uc.Call(id, usecase.UserSelectionType(action)); err != nil {
		h.logger.Error("moderation failed", "err", err)

		h.answerCallback(cb.ID, "ERROR")
		return
	}

	h.answerCallback(cb.ID, "OK")
}

func (h *ModerationHandler) answerCallback(id, text string) {
	if _, err := h.bot.Request(tgbotapi.NewCallback(id, text)); err != nil {
		h.logger.Error("callback failed", "err", err)
	}
}
