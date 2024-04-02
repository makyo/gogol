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
BenchmarkEvolveNaive2d-8            	    1609	   8735432 ns/op	  531148 B/op	     257 allocs/op
BenchmarkEvolveNaive1d-8            	    2317	   5157689 ns/op	  524515 B/op	       1 allocs/op
BenchmarkEvolveScholes-8            	    6955	   1928381 ns/op	 5767259 B/op	      11 allocs/op
BenchmarkEvolveAbrashStruct-8       	   23722	    520863 ns/op	 1055151 B/op	     257 allocs/op
BenchmarkEvolveAbrash-8             	   67591	    167219 ns/op	   72065 B/op	     257 allocs/op
BenchmarkEvolveAbrash1d-8           	  115836	    101178 ns/op	   65536 B/op	       1 allocs/op
BenchmarkEvolveAbrashChangelist-8   	  133270	     82912 ns/op	  106091 B/op	      14 allocs/op
BenchmarkEvolvePrestafford1-8       	  382392	     32084 ns/op	   54899 B/op	      24 allocs/op
PASS
ok  	github.com/makyo/gogol	117.524s
```
