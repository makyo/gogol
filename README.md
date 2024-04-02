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
BenchmarkEvolveNaive2d-8        	    1436	  10173303 ns/op	  531187 B/op	     257 allocs/op
BenchmarkEvolveNaive1d-8        	    2203	   5615695 ns/op	  524527 B/op	       1 allocs/op
BenchmarkEvolveScholes-8        	    7864	   1926182 ns/op	 5767254 B/op	      11 allocs/op
BenchmarkEvolveAbrashStruct-8   	   21536	    538795 ns/op	 1055156 B/op	     257 allocs/op
BenchmarkEvolveAbrash-8         	   72214	    164118 ns/op	   72065 B/op	     257 allocs/op
BenchmarkEvolveAbrash1d-8       	   94778	    110808 ns/op	   65536 B/op	       1 allocs/op
PASS
ok  	github.com/makyo/gogol	86.224s
```
