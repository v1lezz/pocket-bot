package main

import (
	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/v1lezz/pocket-bot/pkg/repository"
	"github.com/v1lezz/pocket-bot/pkg/repository/boltdb"
	"github.com/v1lezz/pocket-bot/pkg/server"
	"github.com/v1lezz/pocket-bot/pkg/telegram"
	pocket "github.com/zhashkevych/go-pocket-sdk"

	"log"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("6060570164:AAGcyVFegRzr-vsm7RS1cdgWxyZfWIRvPuA")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	pocketClient, err := pocket.NewClient("107875-f2350fdc9cf13f5499b9640")
	if err != nil {
		log.Fatal(err)
	}

	db, err := initDB()

	if err != nil {
		log.Fatal(err)
	}

	tokenRepository := boltdb.NewTokenRepository(db)

	telegramBot := telegram.NewBot(bot, pocketClient, tokenRepository, "http://localhost/")
	authorizationServer := server.NewAuthorizationServer(pocketClient, tokenRepository, "https://t.me/v1lezz_pocket_bot")
	go func() {
		if err = telegramBot.Start(); err != nil {
			log.Fatal(err)
		}
	}()
	if err = authorizationServer.Start(); err != nil {
		log.Fatal(err)
	}
}

func initDB() (*bolt.DB, error) {
	db, err := bolt.Open("bot.db", 0600, nil)

	if err != nil {
		log.Fatal(err)
	}

	if err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(repository.AccessTokens))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(repository.RequestTokens))
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return db, nil
}
