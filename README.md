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
BenchmarkEvolveNaive2d-8         	    1458	   8926277 ns/op	  530798 B/op	     257 allocs/op
BenchmarkEvolveNaive1d-8         	    2181	   5333053 ns/op	  524529 B/op	       1 allocs/op
BenchmarkEvolveScholes-8         	    6126	   2015931 ns/op	 5767270 B/op	      11 allocs/op
BenchmarkEvolveAbrash-8          	   19974	    588416 ns/op	 1054775 B/op	     257 allocs/op
BenchmarkEvolveAbrashBitwise-8   	   70332	    165391 ns/op	   71681 B/op	     257 allocs/op
PASS
ok  	github.com/makyo/gogol	69.763s
```
