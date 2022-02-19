package data

import (
	"encoding/json"
	"sync"
)

func (i InputURL) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}

type InputURL struct {
	Id   string
	Urls []string `json:"urls" binding:"required"`
}

type ShortlyURLS struct {
	Parent    InputURL
	Redirects []string
}

func CreateShortlyURL(urls InputURL) ShortlyURLS {
	return ShortlyURLS{Parent: urls, Redirects: make([]string, len(urls.Urls))}
}

func CreateRetrievalURL(urls InputURL) ShortlyURLS {
	return ShortlyURLS{Redirects: urls.Urls, Parent: InputURL{Urls: make([]string, len(urls.Urls))}}
}

type autoInc struct {
	sync.Mutex // ensures autoInc is goroutine-safe
	id         int64
}

var runTimeID autoInc

func (a *autoInc) ID() (id int64) {
	a.Lock()
	defer a.Unlock()

	id = a.id
	a.id++
	return
}
