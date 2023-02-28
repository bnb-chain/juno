package types

// HeightQueue is a simple type alias for a (buffered) channel of block heights.
type HeightQueue chan uint64

func NewQueue(size int) HeightQueue {
	return make(chan uint64, size)
}
