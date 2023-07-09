package boltdb

import (
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
	"github.com/v1lezz/pocket-bot/pkg/repository"
	"strconv"
)

type TokenRepository struct {
	db *bolt.DB
}

func NewTokenRepository(db *bolt.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) Save(chatID int64, token string, bucket repository.Bucket) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Put(int64ToBytes(chatID), []byte(token))
	})

}

func (r *TokenRepository) Get(chatID int64, bucket repository.Bucket) (string, error) {
	var token string
	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		token = string(b.Get(int64ToBytes(chatID)))
		return nil
	})
	if err != nil {
		return "", err
	}

	if token == "" {
		return token, errors.New("token not found")
	}
	return token, nil
}

func int64ToBytes(value int64) []byte {
	return []byte(strconv.FormatInt(value, 10))
}
