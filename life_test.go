package main

import (
	"testing"

	"github.com/makyo/gogol/base"
	"github.com/makyo/gogol/naive1d"
	"github.com/makyo/gogol/naive2d"
	"github.com/makyo/gogol/scholes"
)

var num = 1

func BenchmarkEvolveNaive2d(b *testing.B) {
	var m base.Model
	m = naive2d.New(100, 100, true)
	m = m.Populate()
	for i := 0; i < b.N; i++ {
		m = m.Next()
	}
}

func BenchmarkEvolveNaive1d(b *testing.B) {
	var m base.Model
	m = naive1d.New(100, 100, true)
	m = m.Populate()
	for i := 0; i < b.N; i++ {
		m = m.Next()
	}
}

func BenchmarkEvolveScholes(b *testing.B) {
	var m base.Model
	m = scholes.New(100, 100, true)
	m = m.Populate()
	for i := 0; i < b.N; i++ {
		m = m.Next()
	}
}
