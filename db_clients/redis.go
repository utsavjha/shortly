package db_clients

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	DM "shortly.data.data_model"
	"time"
)

var EXPIRATION_TIME = time.Hour * 4

func StoreURL(conn *redis.Conn, ctx *context.Context, shortlyURL *DM.ShortlyURLS) {
	go func() {
		err := conn.Set(*ctx, shortlyURL.Redirect, shortlyURL.Parent.URL, EXPIRATION_TIME).Err()
		if err != nil {
			panic(fmt.Sprintf("Failed saving key url | Error: %v - shortUrl: %s - originalUrl: %s\n", err, shortlyURL.Redirect, shortlyURL.Parent))
		}
	}()
}

func RetrieveURL(conn *redis.Conn, ctx *context.Context, redisKeyToSearch string) string {
	var parent, err = conn.Get(*ctx, redisKeyToSearch).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed RetrieveInitialUrl url | Error: %v - shortUrl: %s\n", err, redisKeyToSearch))
	}
	return parent
}
