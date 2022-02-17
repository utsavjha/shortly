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

const ServerUrl = "http://shortly:8080/"

func main() {
	var ctx = context.Background()
	var r *gin.Engine = gin.Default()
	//r.SetTrustedProxies([]string{"[::1]"})
	ipChannel := make(chan string, 1000)
	defer close(ipChannel)
	rdb := redis.NewClient(&redis.Options{
		Addr:     "keydb:6379",
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
	//curl -X 'POST' http://shortly:8080/url -H 'content-type: application/json' -d '{"url":"http://abc-xyz.com"}'
	//curl -X 'POST' http://localhost:8080/retrieve -H 'content-type: application/json' -d '{"url":"http://shortly:8080/ddee56b7a0cb1a8af87d4f86ebd7365035e6d29e"}'
	var url DM.URLToShorten
	hn := func(gc *gin.Context) {
		GetIP(gc.Request, IPChannel)
		if err := gc.ShouldBindBodyWith(&url, binding.JSON); err != nil {
			gc.AbortWithError(http.StatusBadRequest, err)
			log.Fatal(http.StatusBadRequest, err)
			return
		}
		url.Id = <-IPChannel
		fmt.Printf("YourIP: %s\n", url.Id)
		shortlyURLS := DM.NewShortlyURLS(url)
		shortlyURLS.Redirect = shortner.ShortenIt(&shortlyURLS.Parent)

		shortlyRedisClient.StoreURL(conn, ctx, shortlyURLS)
		var completeShortURL = fmt.Sprintf("%s%s", ServerUrl, shortlyURLS.Redirect)
		fmt.Printf("Saved shortUrl: %s - originalUrl: %s\n", completeShortURL, shortlyURLS.Parent)
		gc.JSON(http.StatusAccepted, completeShortURL)
	}
	return hn
}

func fetchRedirect(conn *redis.Conn, ctx *context.Context) gin.HandlerFunc {
	//http://localhost:8080/7Wxs5pyEKqqG6gPaJwJSsBndg7MyjXZB74GiNJ3b5Hrn
	//curl -X 'POST' http://localhost:8080/retrieve -H 'content-type: application/json' -d '{"url":"http://localhost:8080/7Wxs5pyEKqqG6gPaJwJSsBndg7MyjXZB74GiNJ3b5Hrn"}'
	var url DM.RetrieveURL
	hn := func(gc *gin.Context) {
		if err := gc.ShouldBindBodyWith(&url, binding.JSON); err != nil {
			gc.JSON(http.StatusBadRequest, err)
			log.Fatal(http.StatusBadRequest, err)
			return
		}
		var redisKey = strings.Split(url.URL, ServerUrl)[1]
		var parent = shortlyRedisClient.RetrieveURL(conn, ctx, redisKey)
		gc.JSON(http.StatusAccepted, fmt.Sprintf("redirecting to --> %s", parent))
	}

	return hn
}

//
//func makeLotsOfRequests() {
//	client := &http.Client{}
//
//	req, _ := http.NewRequest("GET", "http://localhost:8080/url", nil)
//	req.Header.Add("Accept", "application/json")
//	resp, err := client.Do(req)
//
//}

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
