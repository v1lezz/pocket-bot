package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	pocket "github.com/v1lezz/go-pocket-sdk"
	"github.com/v1lezz/pocket-bot/pkg/telegram"
	"log"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("6060570164:AAGcyVFegRzr-vsm7RS1cdgWxyZfWIRvPuA")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	pocketClient, err := pocket.NewClient("108044-21c2d1b481a0a41cb30ec75")
	if err != nil {
		log.Fatal(err)
	}
	telegramBot := telegram.NewBot(bot, pocketClient, "localhost")
	if err := telegramBot.Start(); err != nil {
		log.Fatal(err)
	}
}
