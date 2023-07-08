package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	pocket "github.com/v1lezz/go-pocket-sdk"
	"log"
)

type Bot struct {
	bot          *tgbotapi.BotAPI
	pocketClient *pocket.Client
	redirectUrl  string
}

func NewBot(bot *tgbotapi.BotAPI, client *pocket.Client, redirectUrl string) *Bot {
	return &Bot{bot: bot, pocketClient: client, redirectUrl: redirectUrl}
}

func (b *Bot) Start() error {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	updates, err := b.initUpdatesChannel()
	if err != nil {
		return err
	}

	b.handleUpdates(updates)

	return nil
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
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
		b.handleMessage(update.Message)
	}
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
