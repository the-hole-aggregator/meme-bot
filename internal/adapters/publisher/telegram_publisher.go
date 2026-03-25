package publisher

import (
	"fmt"
	"meme-bot/internal/domain"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TGPublisher struct {
	bot    *tgbotapi.BotAPI
	chatID int64
}

func NewTGPublisher(bot *tgbotapi.BotAPI, chatID int64) *TGPublisher {
	return &TGPublisher{bot: bot, chatID: chatID}
}

func (p *TGPublisher) Publish(meme domain.Meme) error {
	file := tgbotapi.FilePath(fmt.Sprintf("tmp/%s.jpg", meme.SourceID))

	msg := tgbotapi.NewPhoto(p.chatID, file)

	_, err := p.bot.Send(msg)
	return err
}
