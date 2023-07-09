package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/v1lezz/pocket-bot/pkg/repository"
)

func (b *Bot) initAutorizationProcess(message *tgbotapi.Message) error {
	authLink, err := b.generateAuthorizationLink(message.Chat.ID)
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(message.Chat.ID,
		fmt.Sprintf(startReplyTemplate, authLink))
	_, err = b.bot.Send(msg)

	return err
}

func (b *Bot) getAccessToken(chatID int64) (string, error) {
	return b.tokenRepository.Get(chatID, repository.AccessTokens)
}

func (b *Bot) generateAuthorizationLink(chatID int64) (string, error) {
	redirectURL := b.generateRedirectUrl(chatID)
	requestToken, err := b.pocketClient.GetRequestToken(context.Background(), b.redirectUrl)
	if err != nil {
		return "", err
	}

	if err = b.tokenRepository.Save(chatID, requestToken, repository.RequestTokens); err != nil {
		return "", err
	}

	return b.pocketClient.GetAuthorizationURL(requestToken, redirectURL)
}

func (b *Bot) generateRedirectUrl(chatID int64) string {
	return fmt.Sprintf("%s?chat_id=%d", b.redirectUrl, chatID)
}
