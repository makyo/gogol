package base

type Model interface {
	Next() Model
	Populate() Model
	Ingest([][]int) Model
	ToggleCell(int, int) Model
	String() string
}
