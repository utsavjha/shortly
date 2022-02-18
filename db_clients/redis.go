package db_clients

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	DM "shortly.data.data_model"
	"time"
)

const ExpirationTime = time.Hour * 4

func StoreURL(conn *redis.Conn, ctx *context.Context, shortlyURL *DM.ShortlyURLS) {
	go func() {
		for idx, url := range shortlyURL.Redirects {
			err := conn.Set(*ctx, url, shortlyURL.Parent.Urls[idx], ExpirationTime).Err()
			if err != nil {
				log.Println(fmt.Sprintf("Error!!! Failed saving key url | Error: %v - shortUrl: %s - originalUrl: %s\n", err, shortlyURL.Redirects, shortlyURL.Parent))
				panic(fmt.Sprintf("Failed saving key url | Error: %v - shortUrl: %s - originalUrl: %s\n", err, shortlyURL.Redirects, shortlyURL.Parent))
			}
		}
	}()
}

func RetrieveURL(conn *redis.Conn, ctx *context.Context, redisKeyToSearch string,
	redirectChannel chan<- string) {
	var parent, err = conn.Get(*ctx, redisKeyToSearch).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed RetrieveInitialUrl url | Error: %v - shortUrl: %s\n", err, redisKeyToSearch))
	}
	redirectChannel <- parent
}
