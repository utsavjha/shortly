package shortner

import (
	"crypto/sha1"
	"fmt"
)

func ShortenIt(inpURL string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(inpURL)))
}
