package publisher

import (
	"fmt"
	"meme-bot/internal/domain"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ModerationPublisher struct {
	bot    *tgbotapi.BotAPI
	chatID int64
}

func NewModerationPublisher(bot *tgbotapi.BotAPI, chatID int64) *ModerationPublisher {
	return &ModerationPublisher{bot: bot, chatID: chatID}
}

func (p *ModerationPublisher) Publish(meme domain.Meme) error {
	file := tgbotapi.FilePath(fmt.Sprintf("tmp/%s.jpg", meme.SourceID))

	msg := tgbotapi.NewPhoto(p.chatID, file)

	msg.Caption = fmt.Sprintf("ID: %d SourceID: %v", meme.ID, meme.SourceID)

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👍", fmt.Sprintf("approved:%d", meme.ID)),
			tgbotapi.NewInlineKeyboardButtonData("👎", fmt.Sprintf("rejected:%d", meme.ID)),
		),
	)

	if _, err := p.bot.Send(msg); err != nil {
		return err
	}

	return nil
}
