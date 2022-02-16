package data

import (
	"encoding/json"
	"sync"
)

func (i URLToShorten) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}

type URLToShorten struct {
	Id  string
	URL string `form:"inpURL" json:"url" binding:"required"`
}

type RetrieveURL struct {
	URL string `form:"inpURL" json:"url" binding:"required"`
}

type ShortlyURLS struct {
	Parent   URLToShorten
	Redirect string
}

func NewShortlyURLS(url URLToShorten) *ShortlyURLS {
	return &ShortlyURLS{Parent: url}
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
