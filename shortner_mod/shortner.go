package shortner

import (
	"crypto/sha1"
	"fmt"
	DM "shortly.data.data_model"
)

func ShortenIt(shortlyURLS DM.ShortlyURLS) {
	for idx, url := range shortlyURLS.Parent.Urls {
		shortlyURLS.Redirects[idx] = fmt.Sprintf("%x", sha1.Sum([]byte(url+shortlyURLS.Parent.Id)))
	}
}
