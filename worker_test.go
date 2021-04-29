package adj

import (
	"testing"
)

func TestSet(t *testing.T) {
	tests := []struct {
		name     string
		parallel int
		max      int
		want     int
	}{
		{
			name:     "less than",
			parallel: 5,
			max:      10,
			want:     5,
		},
		{
			name:     "equal",
			parallel: 10,
			max:      10,
			want:     10,
		},
		{
			name:     "more than",
			parallel: 15,
			max:      10,
			want:     10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Set(tt.parallel, tt.max)
			if got != tt.want {
				t.Errorf("want %d, got %d", tt.want, got)
			}
		})
	}
}

// TestDo test actual http calls to external resources.
// As the success of test depends on network connection, it is probably not
// a good idea to include this as unit tests.
// A mock HTTP server with mux would be preferable perhaps.
func TestDo(t *testing.T) {
	want := map[string]string{
		"https://xkcd.com/1/info.0.json":  "efd01a3ec9474ef4d2fdde921bde9269",
		"https://xkcd.com/2/info.0.json":  "27d506f481ec17bd177171f282352242",
		"https://xkcd.com/3/info.0.json":  "bcc605204d79e01d61be6a4690ec7539",
		"https://xkcd.com/4/info.0.json":  "52d99ff0a2035ab3b2b3187c3835f79f",
		"https://xkcd.com/5/info.0.json":  "dc1a1ef43b5a365570df0daabad58f7e",
		"https://xkcd.com/6/info.0.json":  "00e43fce497bafaf753f0b7fa3fe4a33",
		"https://xkcd.com/7/info.0.json":  "fbf1f5a92cd5b71ff6cc537e16565fbb",
		"https://xkcd.com/8/info.0.json":  "f5da3b73b71e944828dd1614f17f7b13",
		"https://xkcd.com/9/info.0.json":  "76a5365f33e95a6f7e358524cbc35710",
		"https://xkcd.com/10/info.0.json": "825904392c4ba4a063d4236c43af2679",
	}
	var urls []string
	for url := range want {
		urls = append(urls, url)
	}
	parallel := 3
	jobs := make(chan string, len(urls))
	results := make(chan *Hashed, parallel)
	errs := make(chan error, parallel)

	for i := 0; i < parallel; i++ {
		go Do(jobs, results, errs)
	}

	for _, url := range urls {
		jobs <- url
	}
	close(jobs)

	for i := 0; i < len(urls); i++ {
		if err := <-errs; err != nil {
			continue
		}
		got := <-results
		if got.Response != want[got.URL] {
			t.Errorf("want %s, got %s", want[got.URL], got.Response)
		}
	}
}
