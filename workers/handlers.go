package workers

import (
	"context"
	"github.com/go-redis/redis/v8"
	DM "shortly.data.data_model"
	shortlyRedisClient "shortly.db.clients"
	"shortner"
	"strings"
)

const ShortlyServerURL = "http://shortly:8080/"

func ShortenInputURLs(conn *redis.Conn, ctx *context.Context,
	shortlyChannel <-chan *DM.ShortlyURLS, ipChannel <-chan string, shortenedURLChannel chan<- *DM.ShortlyURLS) {
	for ip := range ipChannel {
		go func(ipAddress string) {
			for url := range shortlyChannel {
				url.Parent.Id = ipAddress
				shortner.ShortenIt(url)
				shortlyRedisClient.StoreURL(conn, ctx, url)
				shortenedURLChannel <- url
			}
		}(ip)
	}
}

func RetrieveParentsOfShortlyURLs(conn *redis.Conn, ctx *context.Context,
	inputShortlyURLChannel <-chan DM.RetrieveURL, parents chan<- string) {
	for shortlyURL := range inputShortlyURLChannel {
		go func(url string) {
			shortlyRedisClient.RetrieveURL(conn, ctx, strings.Split(url, ShortlyServerURL)[1], parents)
		}(shortlyURL.URL)
	}

}
