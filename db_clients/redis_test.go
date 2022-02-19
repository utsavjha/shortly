package db_clients

import (
	"context"
	"github.com/stretchr/testify/assert"
	"reflect"
	data "shortly.data.data_model"
	"testing"
	"time"
)

var ctx = context.Background()
var rdb = CreateConnection()
var conn = rdb.Conn(ctx)
var expiry = time.Second * 2
var computedHash = "HashedValue"
var testURL = "http://abc-xyz.com"
var shortlyURL = data.CreateShortlyURL(data.InputURL{Urls: []string{testURL}})

func TestStoreURL(t *testing.T) {
	shortlyURL.Redirects = []string{computedHash}
	StoreURL(conn, &ctx, shortlyURL, expiry)
	var observedUrl, err = conn.Get(ctx, computedHash).Result()
	assert.Equal(t, testURL, observedUrl,
		testURL, observedUrl, err)

	// now test the expiry
	time.Sleep(expiry * 2)
	var someUrl, err2 = conn.Get(ctx, computedHash).Result()
	assert.Empty(t, someUrl, "Expected Empty but got "+someUrl)
	assert.NotNil(t, err2, observedUrl)
}

func TestFetchKey(t *testing.T) {
	shortlyURL.Redirects = []string{computedHash}
	urlChannel := make(chan string, 1)
	expectedChannel := make(chan string, 1)
	defer close(urlChannel)
	defer close(expectedChannel)
	expectedChannel <- testURL
	conn.Set(ctx, computedHash, testURL, expiry*5).Err()
	FetchKey(conn, &ctx, &computedHash, urlChannel)

	observedValue := <-urlChannel
	assert.Equal(t, testURL, observedValue, testURL, observedValue)
	reflect.DeepEqual(expectedChannel, urlChannel)
}
