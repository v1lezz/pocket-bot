package main

import (
	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/v1lezz/pocket-bot/pkg/config"
	"github.com/v1lezz/pocket-bot/pkg/repository"
	"github.com/v1lezz/pocket-bot/pkg/repository/boltdb"
	"github.com/v1lezz/pocket-bot/pkg/server"
	"github.com/v1lezz/pocket-bot/pkg/telegram"
	pocket "github.com/zhashkevych/go-pocket-sdk"

	"log"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	pocketClient, err := pocket.NewClient(cfg.PocketConsumerKey)
	if err != nil {
		log.Fatal(err)
	}

	db, err := initDB(cfg)

	if err != nil {
		log.Fatal(err)
	}

	tokenRepository := boltdb.NewTokenRepository(db)

	telegramBot := telegram.NewBot(bot, pocketClient, tokenRepository, cfg.AuthServerURL, cfg.Messages)
	authorizationServer := server.NewAuthorizationServer(pocketClient, tokenRepository, cfg.TelegramBotURL)
	go func() {
		if err = telegramBot.Start(); err != nil {
			log.Fatal(err)
		}
	}()
	if err = authorizationServer.Start(); err != nil {
		log.Fatal(err)
	}
}

func initDB(cfg *config.Config) (*bolt.DB, error) {
	db, err := bolt.Open(cfg.DBPath, 0600, nil)

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
