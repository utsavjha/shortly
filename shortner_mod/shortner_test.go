package shortner

import (
	"github.com/stretchr/testify/assert"
	data "shortly.data.data_model"
	"testing"
)

func TestShortenIt(t *testing.T) {
	var computedHashes = []string{"30c3c257ec2a128fca415ac3f151fa04babde86a", "5c587fe402a4df768f51141c76728e25ffaab1d7"}
	var testURLs = []string{"http://xyz.com", "http://abc.com"}
	var shortlyURL = data.CreateShortlyURL(data.InputURL{Urls: testURLs, Id: "someID"})
	ShortenIt(shortlyURL)

	assert.EqualValues(t, shortlyURL.Redirects, computedHashes)
}
