package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

var (
	errInvalidURL   = errors.New("url is invalid")
	errUnautorized  = errors.New("user is not authorized")
	errUnableToSave = errors.New("unable to save")
)

func (b *Bot) handleError(chatID int64, err error) {
	msg := tgbotapi.NewMessage(chatID, b.messages.Errors.Default)
	switch err {
	case errInvalidURL:
		msg.Text = b.messages.Errors.InvalidURL
	case errUnautorized:
		msg.Text = b.messages.Errors.Unauthorized
	case errUnableToSave:
		msg.Text = b.messages.Errors.UnableToSave
	default:

	}
	b.bot.Send(msg)

}

//msg.Text = "Ты не авторизирован! Используй команду /start"
//_, err = b.bot.Send(msg)
//return err
//msg.Text = "Это невалидная ссылка"
///msg.Text = "Увы, не удалось сохранить ссылку. Попробуй еще раз позже."
