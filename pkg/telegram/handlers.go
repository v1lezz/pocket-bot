package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zhashkevych/go-pocket-sdk"
	"net/url"
)

const (
	commandStart = "start"

	startReplyTemplate     = "Привет! Чтобы сохранить ссылки в своем Pocket аккаунте, для начала тебе необходимо дать мне на это доступ. Для этого перейди по ссылке:\n%s"
	replyAlreadyAuthorized = "Ты уже авторизирован. Присылай ссылку, я ее сохраню."
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)
	default:
		return b.handleUnknownCommand(message)
	}

}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Ссылка успешно сохранена!")
	_, err := url.ParseRequestURI(message.Text)
	if err != nil {
		msg.Text = "Это невалидная ссылка"
	}
	msg.ReplyToMessageID = message.MessageID

	accessToken, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		msg.Text = "Ты не авторизирован! Используй команду /start"
		_, err = b.bot.Send(msg)
		return err
	}
	if err = b.pocketClient.Add(context.Background(), pocket.AddInput{
		AccessToken: accessToken,
		URL:         message.Text,
	}); err != nil {
		msg.Text = "Увы, не удалось сохранить ссылку. Попробуй еще раз позже."
	}
	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	_, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		return b.initAutorizationProcess(message)
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, replyAlreadyAuthorized)
	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Я не знаю такой команды :(")
	_, err := b.bot.Send(msg)
	return err
}
