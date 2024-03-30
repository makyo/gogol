# gogol
Golang Game of Life explorations

## Benchmarking

While this project originally started out as a way to teach myself some [Charm.sh](https://charm.sh) tools, it turned into an algorithm exploration. I am primarily working through Eric Lippert's [series on the topic](https://conwaylife.com/wiki/Tutorials/Coding_Life_simulators).

Benchmarks are run with:

    go test -bench . -benchtime=10s -benchmem

```
goos: linux
goarch: amd64
pkg: github.com/makyo/gogol
cpu: Intel(R) Core(TM) i7-10510U CPU @ 1.80GHz
BenchmarkEvolveNaive2d-8   	    1455	   9156673 ns/op	  530802 B/op	     257 allocs/op
BenchmarkEvolveNaive1d-8   	    2196	   5513237 ns/op	  524529 B/op	       1 allocs/op
BenchmarkEvolveScholes-8   	    6386	   2092751 ns/op	 5767266 B/op	      11 allocs/op
BenchmarkEvolveAbrash-8    	   20288	    594999 ns/op	 1054774 B/op	     257 allocs/op
PASS
ok  	github.com/makyo/gogol	58.424s
```
