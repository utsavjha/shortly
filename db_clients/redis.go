package db_clients

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	DM "shortly.data.data_model"
	"time"
)

const KeyDBURL = "localhost:6379"

func CreateConnection() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     KeyDBURL,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func StoreURL(conn *redis.Conn, ctx *context.Context, shortlyURL DM.ShortlyURLS, expiryTime time.Duration) {
	go func() {
		for idx, url := range shortlyURL.Redirects {
			err := conn.Set(*ctx, url, shortlyURL.Parent.Urls[idx], expiryTime).Err()
			if err != nil {
				log.Println(fmt.Sprintf("Error!!! Failed saving key url | Error: %v - shortUrl: %s - originalUrl: %s\n", err, shortlyURL.Redirects, shortlyURL.Parent))
				panic(fmt.Sprintf("Failed saving key url | Error: %v - shortUrl: %s - originalUrl: %s\n", err, shortlyURL.Redirects, shortlyURL.Parent))
			}
		}
	}()
}

func FetchKey(conn *redis.Conn, ctx *context.Context, redisKeyToSearch *string,
	redirectChannel chan<- string) {
	go func(redisKey *string) {
		var parent, err = conn.Get(*ctx, *redisKey).Result()
		if err != nil {
			panic(fmt.Sprintf("Failed RetrieveInitialUrl url | Error: %v - shortUrl: %s\n", err, *redisKey))
		}
		redirectChannel <- parent
	}(redisKeyToSearch)
}
