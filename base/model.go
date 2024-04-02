package base

import "github.com/makyo/gogol/rle"

type Model interface {
	Next()
	Populate()
	ToggleCell(int, int)
	Ingest(*rle.RLEField)
	String() string
}
