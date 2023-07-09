package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/v1lezz/pocket-bot/pkg/repository"
	pocket "github.com/zhashkevych/go-pocket-sdk"
	"log"
)

type Bot struct {
	bot             *tgbotapi.BotAPI
	pocketClient    *pocket.Client
	tokenRepository repository.TokenRepository
	redirectUrl     string
}

func NewBot(bot *tgbotapi.BotAPI, client *pocket.Client, tr repository.TokenRepository, redirectUrl string) *Bot {
	return &Bot{bot: bot, pocketClient: client, tokenRepository: tr, redirectUrl: redirectUrl}
}

func (b *Bot) Start() error {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	updates, err := b.initUpdatesChannel()
	if err != nil {
		return err
	}

	err = b.handleUpdates(updates)

	return err
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) error {
	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.IsCommand() {
			err := b.handleCommand(update.Message)
			if err != nil {
				log.Fatal(err)
			}
			continue
		}
		err := b.handleMessage(update.Message)
		return err
	}
	return nil
}

func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.bot.GetUpdatesChan(u)
	if err != nil {
		return nil, err
	}
	return updates, nil
}
