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
	userRetrievalInputsChannel := make(chan string, 1000)

	defer close(ipChannel)
	defer close(userLongURLInputChannel)
	defer close(userShortlyOutputChannel)
	defer close(shortenedURLChannel)
	defer close(userRetrievalInputsChannel)
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
	r.POST("/retrieve", fetchRedirect(c2, &ctx, userRetrievalInputsChannel, shortenedURLChannel))
	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	r.Run()
	c1.Close()
	c2.Close()
}

func shortenURL(conn *redis.Conn, ctx *context.Context, IPChannel chan string, userInputsChannel chan DM.ShortlyURLS, shortlyURLOutputChannel chan DM.ShortlyURLS) gin.HandlerFunc {
	//curl -X 'POST' http://shortly:8080/url -H 'content-type: application/json' -d '{"urls": ["http://abcxyz.com","http://facebook.com","http://google.com","http://amazon.com","http://bol.com"]}'
	//curl -X 'POST' http://shortly:8080/url -H 'content-type: application/json' -d '{"urls": [,"http://shortly:8080/a9bb4496370625c19838c7f812a1c0e946938dec" ,"http://shortly:8080/8f2453e6d59f946e89bd77bce564cc3652a58843" ,"http://shortly:8080/7d9f902afa6273f852bfa737e1ab6205e3e06ebd", "http://shortly:8080/9308539e8d36064043c4089803d9610c31ab1ee7"]}'
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

func fetchRedirect(conn *redis.Conn, ctx *context.Context, URLRetrievalChannel chan string, outputRedirectsChannel chan string) gin.HandlerFunc {
	//curl -X 'POST' http://shortly:8080/retrieve -H 'content-type: application/json' -d '{"url":"a9bb4496370625c19838c7f812a1c0e946938dec"}'
	//8f2453e6d59f946e89bd77bce564cc3652a58843 7d9f902afa6273f852bfa737e1ab6205e3e06ebd 9308539e8d36064043c4089803d9610c31ab1ee7 ea3e27b2ddf6534331df3325395b18836c0a417e]
	var url DM.InputURL
	hn := func(gc *gin.Context) {
		if err := gc.ShouldBindBodyWith(&url, binding.JSON); err != nil {
			gc.JSON(http.StatusBadRequest, err)
			log.Fatal(http.StatusBadRequest, err)
			return
		}
		redirectURL := worker.RetrieveParentsOfShortlyURLs(conn, ctx, url, URLRetrievalChannel, outputRedirectsChannel)
		gc.JSON(http.StatusOK, fmt.Sprintf("redirecting %v to ---- %v", redirectURL.Redirects, redirectURL.Parent.Urls))
	}
	return hn
}
