package main

import (
	"fmt"
	"time"

	"github.com/makyo/gogol/abrash"
	"github.com/makyo/gogol/abrash1d"
	"github.com/makyo/gogol/abrashchangelist"
	"github.com/makyo/gogol/abrashstruct"
	"github.com/makyo/gogol/base"
	"github.com/makyo/gogol/naive1d"
	"github.com/makyo/gogol/naive2d"
	"github.com/makyo/gogol/prestafford1"
	"github.com/makyo/gogol/prestafford2"
	"github.com/makyo/gogol/rle"
	"github.com/makyo/gogol/scholes"
)

func main() {
	f, err := rle.Unmarshal(`#N Acorn
#O Charles Corderman
#C A methuselah with lifespan 5206.
#C www.conwaylife.com/wiki/index.php?title=Acorn
x = 7, y = 3, rule = B3/S23
bo5b$3bo3b$2o2b3o!`)
	if err != nil {
		panic(err)
	}
	table := map[string]base.Model{
		"naive2d":          naive2d.New(256, 256),
		"naive1d":          naive1d.New(256, 256),
		"scholes":          scholes.New(256, 256),
		"abrashstruct":     abrashstruct.New(256, 256),
		"abrash":           abrash.New(256, 256),
		"abrash1d":         abrash1d.New(256, 256),
		"abrashchangelist": abrashchangelist.New(256, 256),
		"prestafford1":     prestafford1.New(256, 256),
		"prestafford2":     prestafford2.New(256, 256),
	}
	result := map[string]time.Duration{}

	for k, _ := range table {
		table[k].Ingest(f)
		start := time.Now()
		for i := 0; i < 5000; i++ {
			table[k].Next()
		}
		result[k] = time.Now().Sub(start)
	}
	fmt.Printf("%-20s  time(ms)\n", "Algorithm")
	for k, _ := range result {
		fmt.Printf("%20s %6vms\n", k, result[k].Milliseconds())
	}
}
