package prestafford2

// A triplet of cell is represented by a set of 16 bits, which holds the three cells' current states, next states, and count of neighbors.
type cell uint16

const (
	// The next state of the cell is stored in the 13th, 14th, and 15th bits, 1-indexed
	leftNext   = 14
	middleNext = 13
	rightNext  = 12

	// The current state of the cell is stored in the 10th, 11th, and 12th bits.
	leftState   = 11
	middleState = 10
	rightState  = 9

	// The count of neighbors for each cell in the triplet takes three bits to represent. For the left cell, it is the seventh through ninth, for the middle the fourth through sixth, and for the right the first through third
	leftCount   = 6
	middleCount = 3
	rightCount  = 0

	// Since we most often just add one to the neighbor count, these are the particular "ones" we will be adding, since the value of one is different per slot in the triplet.
	leftCountone   = 1 << leftCount
	middleCountone = 1 << middleCount
	rightCountone  = 1 << rightCount

	// The actual masks used for interacting with the next state, current state, and counts for each of the triplets.
	leftNextbit   = 1 << leftNext
	middleNextbit = 1 << middleNext
	rightNextbit  = 1 << rightNext

	leftStatebit   = 1 << leftState
	middleStatebit = 1 << middleState
	rightStatebit  = 1 << rightState

	leftCountbit   = 7 << leftCount
	middleCountbit = 7 << middleCount
	rightCountbit  = 7 << rightCount

	// This bit will be used for getting the current state of the entire triplet to check against the *next* state of the entire triplet.
	tripletStatebit = leftStatebit | middleStatebit | rightStatebit

	// The next state of the entire triplet.
	tripletNextbit = leftNextbit | middleNextbit | rightNextbit
)

// Here follows a set of utility functions that just surface more data on the cells.

// These functions get the next states of each of the three cells represented in the uint16 by selecting the appropriate bits. The *Raw functions also return that state as a 1 or 0 so we can do math on them.
func (c cell) leftNextRaw() uint16   { return uint16((c & leftNextbit) >> leftNext) }
func (c cell) middleNextRaw() uint16 { return uint16((c & middleNextbit) >> middleNext) }
func (c cell) rightNextRaw() uint16  { return uint16((c & rightNextbit) >> rightNext) }

func (c cell) leftNext() bool   { return (c & leftNextbit) != 0 }
func (c cell) middleNext() bool { return (c & middleNextbit) != 0 }
func (c cell) rightNext() bool  { return (c & rightNextbit) != 0 }

func (c cell) tripletNext() uint16 { return uint16((tripletNextbit & c) >> rightNext) }

// These functions get the current states of each of the three cells represented in the uint16 by selecting the appropriate bits. The *Raw functions also return that state as a 1 or 0 so we can do math on them.
func (c cell) leftStateRaw() uint16   { return uint16((c & leftStatebit) >> leftState) }
func (c cell) middleStateRaw() uint16 { return uint16((c & middleStatebit) >> middleState) }
func (c cell) rightStateRaw() uint16  { return uint16((c & rightStatebit) >> rightState) }

func (c cell) leftState() bool   { return (c & leftStatebit) != 0 }
func (c cell) middleState() bool { return (c & middleStatebit) != 0 }
func (c cell) rightState() bool  { return (c & rightStatebit) != 0 }

func (c cell) tripletState() uint16 { return uint16((tripletStatebit & c) >> rightState) }

// changed lets us know if any of the cells in the triplet will change on this step — this way, we can skip the entire triplet if not.
func (c cell) changed() bool { return c.tripletNext() != c.tripletState() }

// These functions get the current states of each of the three cells represented in the uint16 by selecting the appropriate bits. The *Raw functions return just the counts contained in their bits, while the non *Raw types add the states of the bits to either side (e.g: the left bit will have in its neighbors SE, S, SW, W, NW, N, NE, but since the middle bit is its E neighbor, that needs to be added in).
func (c cell) leftNeighborsRaw() uint16   { return uint16((leftCountbit & c) >> leftCount) }
func (c cell) middleNeighborsRaw() uint16 { return uint16((middleCountbit & c) >> middleCount) }
func (c cell) rightNeighborsRaw() uint16  { return uint16((rightCountbit & c) >> rightCount) }

func (c cell) leftNeighbors() uint16 { return c.middleStateRaw() + c.leftNeighborsRaw() }
func (c cell) middleNeighbors() uint16 {
	return c.leftStateRaw() + c.rightStateRaw() + c.middleNeighborsRaw()
}
func (c cell) rightNeighbors() uint16 { return c.middleStateRaw() + c.rightNeighborsRaw() }

// These functions set the next states of the left, middle, and right cells' bits.
func (c cell) setLeftNext(b bool) cell {
	if b {
		return c | leftNextbit
	} else {
		return c &^ leftNextbit
	}
}
func (c cell) setMiddleNext(b bool) cell {
	if b {
		return c | middleNextbit
	} else {
		return c &^ middleNextbit
	}
}
func (c cell) setRightNext(b bool) cell {
	if b {
		return c | rightNextbit
	} else {
		return c &^ rightNextbit
	}
}

// These functions set the current states of the left, middle, and right cells' bits.
func (c cell) setLeftState(b bool) cell {
	if b {
		return c | leftStatebit
	} else {
		return c &^ leftStatebit
	}
}
func (c cell) setMiddleState(b bool) cell {
	if b {
		return c | middleStatebit
	} else {
		return c &^ middleStatebit
	}
}
func (c cell) setRightState(b bool) cell {
	if b {
		return c | rightStatebit
	} else {
		return c &^ rightStatebit
	}
}
