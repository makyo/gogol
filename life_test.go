package main

import (
	"testing"

	"github.com/makyo/gogol/abrash"
	"github.com/makyo/gogol/abrash1d"
	"github.com/makyo/gogol/abrashchangelist"
	"github.com/makyo/gogol/abrashstruct"
	"github.com/makyo/gogol/naive1d"
	"github.com/makyo/gogol/naive2d"
	"github.com/makyo/gogol/prestafford1"
	"github.com/makyo/gogol/rle"
	"github.com/makyo/gogol/scholes"
)

var (
	num = 1
)

func acorn() *rle.RLEField {
	f, err := rle.Unmarshal(`#N Acorn
#O Charles Corderman
#C A methuselah with lifespan 5206.
#C www.conwaylife.com/wiki/index.php?title=Acorn
x = 7, y = 3, rule = B3/S23
bo5b$3bo3b$2o2b3o!`)
	if err != nil {
		panic(err)
	}
	return f
}

func BenchmarkEvolveNaive2d(b *testing.B) {
	m := naive2d.New(256, 256)
	m.Ingest(acorn())
	for i := 0; i < b.N; i++ {
		m.Next()
	}
}

func BenchmarkEvolveNaive1d(b *testing.B) {
	m := naive1d.New(256, 256)
	m.Ingest(acorn())
	for i := 0; i < b.N; i++ {
		m.Next()
	}
}

func BenchmarkEvolveScholes(b *testing.B) {
	m := scholes.New(256, 256)
	m.Ingest(acorn())
	for i := 0; i < b.N; i++ {
		m.Next()
	}
}

func BenchmarkEvolveAbrashStruct(b *testing.B) {
	m := abrashstruct.New(256, 256)
	m.Ingest(acorn())
	for i := 0; i < b.N; i++ {
		m.Next()
	}
}

func BenchmarkEvolveAbrash(b *testing.B) {
	m := abrash.New(256, 256)
	m.Ingest(acorn())
	for i := 0; i < b.N; i++ {
		m.Next()
	}
}

func BenchmarkEvolveAbrash1d(b *testing.B) {
	m := abrash1d.New(256, 256)
	m.Ingest(acorn())
	for i := 0; i < b.N; i++ {
		m.Next()
	}
}

func BenchmarkEvolveAbrashChangelist(b *testing.B) {
	m := abrashchangelist.New(256, 256)
	m.Ingest(acorn())
	for i := 0; i < b.N; i++ {
		m.Next()
	}
}

func BenchmarkEvolvePrestafford1(b *testing.B) {
	m := prestafford1.New(256, 256)
	m.Ingest(acorn())
	for i := 0; i < b.N; i++ {
		m.Next()
	}
}
