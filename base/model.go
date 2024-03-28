package base

type Model interface {
	Next() Model
	Populate() Model
	ToggleCell(int, int) Model
	String() string
}
