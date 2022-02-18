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
	worker "shortly.workers"
)

const KeyDBURL = "keydb:6379"

func main() {
	var ctx = context.Background()
	var r *gin.Engine = gin.Default()
	//r.SetTrustedProxies([]string{"[::1]"})
	ipChannel := make(chan string, 1000)
	userLongURLInputChannel := make(chan *DM.ShortlyURLS, 1000)
	userShortlyOutputChannel := make(chan *DM.ShortlyURLS, 1000)
	userInputShortlyURLsChannel := make(chan DM.RetrieveURL, 1000)
	shortenedURLChannel := make(chan string, 1000)

	defer close(ipChannel)
	defer close(userLongURLInputChannel)
	defer close(userShortlyOutputChannel)
	defer close(shortenedURLChannel)
	defer close(userInputShortlyURLsChannel)
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
	r.POST("/url", shortenURL(c1, &ctx, ipChannel, userLongURLInputChannel, userShortlyOutputChannel))
	r.POST("/retrieve", fetchRedirect(c2, &ctx, userInputShortlyURLsChannel))
	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	r.Run()
	c1.Close()
	c2.Close()
}

func shortenURL(conn *redis.Conn, ctx *context.Context, IPChannel chan string, userInputsChannel chan *DM.ShortlyURLS, shortlyURLOutputChannel chan *DM.ShortlyURLS) gin.HandlerFunc {
	//curl -X 'POST' http://shortly:8080/url -H 'content-type: application/json' -d '{"urls": ["http://abcxyz.com","http://facebook.com","http://google.com","http://amazon.com","http://bol.com"]}'
	var url DM.URLToShorten
	hn := func(gc *gin.Context) {
		GetIP(gc.Request, IPChannel)
		if err := gc.ShouldBindBodyWith(&url, binding.JSON); err != nil {
			gc.AbortWithError(http.StatusBadRequest, err)
			log.Fatal(http.StatusBadRequest, err)
			return
		}
		userInputsChannel <- DM.CreateShortlyURL(url)
		worker.ShortenInputURLs(conn, ctx, userInputsChannel, IPChannel, shortlyURLOutputChannel)
		for shortlyURL := range shortlyURLOutputChannel {
			gc.JSON(http.StatusOK, fmt.Sprintf("%s --> %s", shortlyURL.Parent.Urls, shortlyURL.Redirects))
		}
	}
	return hn
}

func fetchRedirect(conn *redis.Conn, ctx *context.Context, inputShortlyURLsChannel chan DM.RetrieveURL) gin.HandlerFunc {
	//curl -X 'POST' http://shortly:8080/retrieve -H 'content-type: application/json' -d '{"url":"http://shortly:8080/2ccf2e089861edd16a48f5b83e91e9b3cb4852be"}'
	var url DM.RetrieveURL
	parentsChannel := make(chan string)
	defer close(parentsChannel)
	hn := func(gc *gin.Context) {
		if err := gc.ShouldBindBodyWith(&url, binding.JSON); err != nil {
			gc.JSON(http.StatusBadRequest, err)
			log.Fatal(http.StatusBadRequest, err)
			return
		}
		inputShortlyURLsChannel <- url
		worker.RetrieveParentsOfShortlyURLs(conn, ctx, inputShortlyURLsChannel, parentsChannel)
		for shortenedURL := range parentsChannel {
			gc.JSON(http.StatusOK, fmt.Sprintf("redirecting to --> %s", shortenedURL))
		}
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
