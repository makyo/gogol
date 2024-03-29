package base

type Model interface {
	Next()
	Populate()
	Ingest([][]int)
	ToggleCell(int, int)
	String() string
}
