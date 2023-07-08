package telegram

import (
	"context"
	"fmt"
)

func (b *Bot) generateAuthorizationLink(chatID int64) (string, error) {
	redirectURL := b.generateRedirectUrl(chatID)
	requestToken, err := b.pocketClient.GetRequestToken(context.Background(), b.redirectUrl)
	if err != nil {
		return "", err
	}
	return b.pocketClient.GetAuthorizationURL(requestToken, redirectURL)
}

func (b *Bot) generateRedirectUrl(chatID int64) string {
	return fmt.Sprintf("%s?chat_id=%d", b.redirectUrl, chatID)
}
