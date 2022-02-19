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

const KeyDBURL = "localhost:6379"

func main() {
	var ctx = context.Background()
	var r *gin.Engine = gin.Default()
	//r.SetTrustedProxies([]string{"[::1]"})
	ipChannel := make(chan string, 1000)
	userLongURLInputChannel := make(chan DM.ShortlyURLS, 1000)
	userShortlyOutputChannel := make(chan DM.ShortlyURLS, 1000)
	shortenedURLChannel := make(chan string, 1000)

	defer close(ipChannel)
	defer close(userLongURLInputChannel)
	defer close(userShortlyOutputChannel)
	defer close(shortenedURLChannel)
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
	r.POST("/shorten", shortenURL(c1, &ctx, ipChannel, userLongURLInputChannel, userShortlyOutputChannel))
	r.POST("/retrieve", fetchRedirect(c2, &ctx, shortenedURLChannel))
	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	r.Run()
	c1.Close()
	c2.Close()
}

func shortenURL(conn *redis.Conn, ctx *context.Context, IPChannel chan string, userInputsChannel chan DM.ShortlyURLS, shortlyURLOutputChannel chan DM.ShortlyURLS) gin.HandlerFunc {
	var url DM.InputURL
	hn := func(gc *gin.Context) {
		worker.GetIP(gc.Request, IPChannel)
		if err := gc.ShouldBindBodyWith(&url, binding.JSON); err != nil {
			gc.AbortWithError(http.StatusBadRequest, err)
			log.Fatal(http.StatusBadRequest, err)
			return
		}
		shortlyURL := worker.GetShortlyURL(conn, ctx, url, IPChannel, userInputsChannel, shortlyURLOutputChannel)
		gc.JSON(http.StatusOK, fmt.Sprintf("%v -- %v", shortlyURL.Parent.Urls, shortlyURL.Redirects))
	}
	return hn
}

func fetchRedirect(conn *redis.Conn, ctx *context.Context, outputRedirectsChannel chan string) gin.HandlerFunc {
	var url DM.InputURL
	hn := func(gc *gin.Context) {
		if err := gc.ShouldBindBodyWith(&url, binding.JSON); err != nil {
			gc.JSON(http.StatusBadRequest, err)
			log.Fatal(http.StatusBadRequest, err)
			return
		}
		redirectURL := worker.RetrieveParentsOfShortlyURLs(conn, ctx, url, outputRedirectsChannel)
		gc.JSON(http.StatusOK, fmt.Sprintf("redirecting %v to ---- %v", redirectURL.Redirects, redirectURL.Parent.Urls))
	}
	return hn
}
