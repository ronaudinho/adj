package adj

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Hashed contains modified URL and its MD5 hashed response.
type Hashed struct {
	URL      string
	Response string
}

// Get fetch from a given URL, adding protocol to it if necessary.
// The protocol defaults to http.
// It returns the updated URL with MD5 hash of response body
// or an error depending on the result of the call.
// This function is synchronous by design and
// it is the caller's choice to concurrently call it.
// https://talks.golang.org/2013/bestpractices.slide#25
func Get(url string) (*Hashed, error) {
	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	h := md5.New()
	// while this works for most cases, it might be necessary to
	// pipe the response body instead for page where response is large
	_, err = io.Copy(h, res.Body)
	if err != nil {
		return nil, err
	}
	return &Hashed{
		URL:      url,
		Response: fmt.Sprintf("%x", h.Sum(nil)),
	}, nil
}
