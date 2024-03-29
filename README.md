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
BenchmarkEvolveNaive2d-8   	    1353	  11982964 ns/op	  531220 B/op	     257 allocs/op
BenchmarkEvolveNaive1d-8   	    1861	   5965848 ns/op	  524857 B/op	       1 allocs/op
BenchmarkEvolveScholes-8   	    5846	   2186120 ns/op	 5767457 B/op	      11 allocs/op
BenchmarkEvolveAbrash-8    	   18666	    598567 ns/op	 1054806 B/op	     257 allocs/op
PASS
ok  	github.com/makyo/gogol	59.537s
```
