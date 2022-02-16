package shortner

import (
	"crypto/sha1"
	"fmt"
	DM "shortly.data.data_model"
)

func ShortenIt(urlToShorten *DM.URLToShorten) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(urlToShorten.URL+urlToShorten.Id)))
}
