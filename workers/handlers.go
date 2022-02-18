package workers

import (
	"context"
	"github.com/go-redis/redis/v8"
	"net/http"
	DM "shortly.data.data_model"
	shortlyRedisClient "shortly.db.clients"
	"shortner"
	"strings"
)

const ShortlyServerURL = "http://shortly:8080/"

func ShortenInputURLs(conn *redis.Conn, ctx *context.Context,
	shortlyChannel <-chan DM.ShortlyURLS, ip string, shortenedURLChannel chan<- DM.ShortlyURLS) {
	go func() {
		for {
			select {
			case url := <-shortlyChannel:
				url.Parent.Id = ip
				shortner.ShortenIt(url)
				shortlyRedisClient.StoreURL(conn, ctx, url)
				shortenedURLChannel <- url
				return
			}
		}
	}()
}
func DetermineRedisKeys(redirects []string, redisFetchChannel chan<- string) {
	go func() {
		for _, redirectURL := range redirects {
			redisKey := strings.Split(redirectURL, ShortlyServerURL)[1]
			redisFetchChannel <- redisKey
		}
	}()
}
func RetrieveParentsOfShortlyURLs(conn *redis.Conn, ctx *context.Context,
	url DM.InputURL, redisFetchChannel chan string, redirectsChannel chan string) DM.ShortlyURLS {
	shortlyURL := DM.CreateRetrievalURL(url)
	DetermineRedisKeys(shortlyURL.Redirects, redisFetchChannel)
	for {
		select {
		case redisKey := <-redisFetchChannel:
			shortlyRedisClient.FetchKey(conn, ctx, &redisKey, redirectsChannel)
		case shortURL := <-redirectsChannel:

		}
	}
	x := 0
	for shortURL := range redirectsChannel {
		shortlyURL.Parent.Urls[x] = shortURL
		x++
	}
	return shortlyURL

}

func GetShortlyURL(conn *redis.Conn, ctx *context.Context, inpURL DM.InputURL, ipChannel chan string, userInputsChannel chan DM.ShortlyURLS, shortenedURLsOutputChannel chan DM.ShortlyURLS) DM.ShortlyURLS {
	for {
		select {
		case ip := <-ipChannel:
			shortlyURL := DM.CreateShortlyURL(inpURL)
			userInputsChannel <- shortlyURL
			ShortenInputURLs(conn, ctx, userInputsChannel, ip, shortenedURLsOutputChannel)

		case shortURL := <-shortenedURLsOutputChannel:
			return shortURL
		}
	}
}

func GetIP(r *http.Request, IPs chan<- string) {
	// GetIP gets a requests IP address by reading off the forwarded-for
	// header (for proxies) and falls back to use the remote address.
	go func() {
		var forwarded = r.Header.Get("X-FORWARDED-FOR")
		if forwarded == "" {
			IPs <- r.Host
		} else {
			IPs <- forwarded
		}
	}()
}
