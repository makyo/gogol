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
BenchmarkEvolveNaive2d-8   	  19618	   686576 ns/op	 530487 B/op	    257 allocs/op
BenchmarkEvolveNaive1d-8   	   1260	  9047516 ns/op	 525125 B/op	      1 allocs/op
BenchmarkEvolveScholes-8   	   5683	  2177095 ns/op	5767465 B/op	     11 allocs/op
PASS
ok  	github.com/makyo/gogol	44.592s
```
