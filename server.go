package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	DM "shortly.data.data_model"
	shortlyRedisClient "shortly.db.clients"
	"shortner"
	"strings"
)

const ShortlyServerURL = "http://shortly:8080/"
const KeyDBURL = "keydb:6379"

func main() {
	var ctx = context.Background()
	var r *gin.Engine = gin.Default()
	//r.SetTrustedProxies([]string{"[::1]"})
	ipChannel := make(chan string, 1000)
	defer close(ipChannel)
	rdb := redis.NewClient(&redis.Options{
		Addr:     KeyDBURL,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//sampleDataRequest := ""
	c1 := rdb.Conn(ctx)
	c2 := rdb.Conn(ctx)
	if err := c1.ClientSetName(ctx, "shortener").Err(); err != nil {
		panic(err)
	}
	r.POST("/url", shortenURL(c1, &ctx, ipChannel))
	r.POST("/retrieve", fetchRedirect(c2, &ctx))
	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	r.Run()
	c1.Close()
	c2.Close()
}

func shortenURL(conn *redis.Conn, ctx *context.Context, IPChannel chan string) gin.HandlerFunc {
	//curl -X 'POST' http://shortly:8080/url -H 'content-type: application/json' -d '{"urls": ["http://abcxyz.com","http://facebook.com","http://google.com","http://amazon.com","http://bol.com"]}'
	var url DM.URLToShorten
	hn := func(gc *gin.Context) {
		GetIP(gc.Request, IPChannel)
		if err := gc.ShouldBindBodyWith(&url, binding.JSON); err != nil {
			gc.AbortWithError(http.StatusBadRequest, err)
			log.Fatal(http.StatusBadRequest, err)
			return
		}
		shortlyURLS := DM.CreateShortlyURL(url)

		url.Id = <-IPChannel
		//fmt.Printf("YourIP: %s\n", url.Id)
		shortner.ShortenIt(shortlyURLS)
		shortlyRedisClient.StoreURL(conn, ctx, shortlyURLS)

		fmt.Printf("Saved shortUrl from User: %s\n", shortlyURLS.Parent.Id)
		for idx, shortenedURL := range shortlyURLS.Redirects {
			gc.JSON(http.StatusOK, fmt.Sprintf("%s --> %s", shortlyURLS.Parent.Urls[idx], shortenedURL))
		}
	}
	return hn
}

func fetchRedirect(conn *redis.Conn, ctx *context.Context) gin.HandlerFunc {
	//curl -X 'POST' http://shortly:8080/retrieve -H 'content-type: application/json' -d '{"url":"http://shortly:8080/2ccf2e089861edd16a48f5b83e91e9b3cb4852be"}'
	var url DM.RetrieveURL
	hn := func(gc *gin.Context) {
		if err := gc.ShouldBindBodyWith(&url, binding.JSON); err != nil {
			gc.JSON(http.StatusBadRequest, err)
			log.Fatal(http.StatusBadRequest, err)
			return
		}
		var redisKey = strings.Split(url.URL, ShortlyServerURL)[1]
		var parent = shortlyRedisClient.RetrieveURL(conn, ctx, redisKey)
		gc.JSON(http.StatusOK, fmt.Sprintf("redirecting to --> %s", parent))
	}

	return hn
}

// GetIP gets a requests IP address by reading off the forwarded-for
// header (for proxies) and falls back to use the remote address.
func GetIP(r *http.Request, IPs chan<- string) {
	go func() {
		var forwarded = r.Header.Get("X-FORWARDED-FOR")
		if forwarded == "" {
			IPs <- r.Host
		} else {
			IPs <- forwarded
		}
	}()
}
