package base

import "github.com/makyo/gogol/rle"

type Model interface {
	Next()
	Populate()
	Ingest(*rle.RLEField)
	ToggleCell(int, int)
	String() string
}
