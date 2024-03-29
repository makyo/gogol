package main

import (
	"testing"

	"github.com/makyo/gogol/abrash"
	"github.com/makyo/gogol/base"
	"github.com/makyo/gogol/naive1d"
	"github.com/makyo/gogol/naive2d"
	"github.com/makyo/gogol/scholes"
)

var num = 1

func acorn() [][]int {
	field := make([][]int, 256)
	for i, _ := range field {
		field[i] = make([]int, 256)
	}
	field[128][128] = 1
	field[128][129] = 1
	field[130][129] = 1
	field[129][131] = 1
	field[128][132] = 1
	field[128][133] = 1
	field[128][134] = 1
	return field
}

func BenchmarkEvolveNaive2d(b *testing.B) {
	var m base.Model
	m = naive2d.New(256, 256)
	m.Ingest(acorn())
	for i := 0; i < b.N; i++ {
		m.Next()
	}
}

func BenchmarkEvolveNaive1d(b *testing.B) {
	var m base.Model
	m = naive1d.New(256, 256)
	m.Ingest(acorn())
	for i := 0; i < b.N; i++ {
		m.Next()
	}
}

func BenchmarkEvolveScholes(b *testing.B) {
	var m base.Model
	m = scholes.New(256, 256)
	m.Ingest(acorn())
	for i := 0; i < b.N; i++ {
		m.Next()
	}
}

func BenchmarkEvolveAbrash(b *testing.B) {
	var m base.Model
	m = abrash.New(256, 256)
	m.Ingest(acorn())
	for i := 0; i < b.N; i++ {
		m.Next()
	}
}
