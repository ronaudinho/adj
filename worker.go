package adj

// Set sets the number of workers.
// It sets the number to max if parallel is more than max.
func Set(parallel, max int) int {
	if parallel > max {
		return max
	}
	return parallel
}

// Do allows calling Get concurrently,
// spawning a worker for each invocation.
func Do(urls <-chan string, results chan<- *Hashed, errs chan<- error) {
	for url := range urls {
		hashed, err := Get(url)
		results <- hashed
		errs <- err
	}
}
