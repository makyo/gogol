# gogol
Golang Game of Life explorations

## Benchmarking

While this project originally started out as a way to teach myself some [Charm.sh](https://charm.sh) tools, it turned into an algorithm exploration. I am primarily working through Eric Lippert's [series on the topic](https://conwaylife.com/wiki/Tutorials/Coding_Life_simulators).

Current status: a naive implementation of the algorithm using a 2-d arrayis the winner, followed by Scholes's algorithm, then a naive 1-d array. However, I have not done much optimization on the first two yet.

```
goos: linux
goarch: amd64
pkg: github.com/makyo/gogol
cpu: Intel(R) Core(TM) i7-10510U CPU @ 1.80GHz
BenchmarkEvolveNaive2d-8   	  179352	    107834 ns/op	   92336 B/op	     102 allocs/op
BenchmarkEvolveNaive1d-8   	    7784	   1519103 ns/op	   81978 B/op	       2 allocs/op
BenchmarkEvolveScholes-8   	   43047	    253757 ns/op	  901170 B/op	      12 allocs/op
PASS
ok  	github.com/makyo/gogol	66.030s
```
