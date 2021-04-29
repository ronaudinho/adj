package main

import (
	"flag"
	"fmt"

	"github.com/ronaudinho/adj"
)

var max, parallel int

func init() {
	flag.IntVar(&parallel, "parallel", 1, "the number of parallel execution")
	flag.IntVar(&max, "max", 10, "the number of maximum parallel execution allowed")
}

func main() {
	flag.Parse()
	workers := adj.Set(parallel, max)

	urls := flag.Args()
	jobs := make(chan string, len(urls))
	results := make(chan *adj.Hashed, workers)
	errs := make(chan error, workers)

	for i := 0; i < workers; i++ {
		go adj.Do(jobs, results, errs)
	}

	for _, url := range urls {
		jobs <- url
	}
	close(jobs)

	for i := 0; i < len(urls); i++ {
		if err := <-errs; err != nil {
			continue
		}
		res := <-results
		fmt.Println(res.URL, res.Response)
	}
}
