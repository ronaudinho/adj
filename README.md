### getting
#### with go get
`go get github.com/ronaudinho/adj`
#### with git
`git clone https://github.com/ronaudinho/adj`

### building 
from the repository root directory
`cd cmd/adj`
`go build`

### running
from `cmd/adj`, after building the binary (assuming the previous steps are followed)
`./adj -parallel 3 google.com facebook.com twitter.com`
for list of accepted flags, run
`./adj -h`
when `-max` is not set, it defaults to 10.
when `-parallel` is not set, it defaults to 1 (that is, operations are run synchronously).

### testing
from the repository root directory, simply run `go test -v ./...`

to get benchmark of parallel executions, run `go test -bench=.`
this is a sample result of running 1000 jobs on local server with different number of workers.
```
goos: linux
goarch: amd64
pkg: github.com/ronaudinho/adj
cpu: Intel(R) Core(TM) i5-6300U CPU @ 2.40GHz
BenchmarkSynchronous-4   	1000000000	         0.07736 ns/op
BenchmarkParallel/1_worker-4         	1000000000	         0.07810 ns/op
BenchmarkParallel/5_workers-4        	1000000000	         0.03565 ns/op
BenchmarkParallel/10_workers-4       	1000000000	         0.03998 ns/op
BenchmarkParallel/50_workers-4       	1000000000	         0.06695 ns/op
BenchmarkParallel/100_workers-4      	1000000000	         0.08202 ns/op
PASS
ok  	github.com/ronaudinho/adj	3.484s
```
since goroutines depends on the number of available cpu, benchmark shows that past a certain number, having more parallel processes actually results in slower operations. the result should vary depending on the machine. actual http calls are not included in the test as they would be affected by network speed as well, and would warrant further investigation if performance is an issue (probably using `pprof` or tracing or even `eBPF`).

### project structure
as it is a rather simple project, I opted for the root directory package codes with the main program inside `cmd` folder.

### assumptions
1. using parts of this project as dependency is not considered as using non-stdlib since they are internal dependencies.
2. the part of response we are interested in is only the response body.
3. error on http call is not explicitly handled (skip on error).

### technical choices
1. Get exposed as synchronous function per https://talks.golang.org/2013/bestpractices.slide#25

### possible improvements
1. using `sync` and/or `sync/atomic` for parallel operations
2. more explicit error handling
3. limits the time required to fetch a resource
