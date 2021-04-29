package adj

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkSynchronous(b *testing.B) {
	jobsCount := 1000
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	}))
	defer server.Close()
	var urls []string
	for i := 0; i < jobsCount; i++ {
		urls = append(urls, server.URL)
	}
	for _, url := range urls {
		Get(url)
	}
}

func BenchmarkParallel(b *testing.B) {
	benches := []struct {
		name     string
		parallel int
	}{
		{"1 worker", 1},
		{"5 workers", 5},
		{"10 workers", 10},
		{"50 workers", 50},
		{"100 workers", 100},
	}

	jobsCount := 1000
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	}))
	defer server.Close()
	var urls []string
	for i := 0; i < jobsCount; i++ {
		urls = append(urls, server.URL)
	}

	for _, bb := range benches {
		b.Run(bb.name, func(b *testing.B) {
			workersCount := bb.parallel
			jobs := make(chan string, jobsCount)
			results := make(chan *Hashed, workersCount)
			errs := make(chan error, workersCount)

			for i := 0; i < workersCount; i++ {
				go Do(jobs, results, errs)
			}

			for _, url := range urls {
				jobs <- url
			}
			close(jobs)

			for i := 0; i < jobsCount; i++ {
				if err := <-errs; err != nil {
					continue
				}
				<-results
			}
		})
	}
}
