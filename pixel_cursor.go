package underbyte

import (
	"math/rand"
)

type Sequence interface {
	next() (int, bool)
	position() int
	limit() int
}

type SequentialSequence struct {
	index int
	max   int
}

type RandomizedSequence struct {
	min           int
	max           int
	maxIterations int

	seen map[int]bool
}

type PixelCursor struct {
	Sequence
}

func (seq *SequentialSequence) position() int {
	return seq.index
}

func (seq *SequentialSequence) limit() int {
	return seq.max
}

func (seq *SequentialSequence) next() (int, bool) {
	if seq.index > seq.max-1 {
		return -1, false
	}

	current := seq.index
	seq.index++
	return current, true
}

func (rseq *RandomizedSequence) position() int {
	return rseq.min
}

func (rseq *RandomizedSequence) limit() int {
	return rseq.max
}

func (rseq *RandomizedSequence) next() (int, bool) {
	if len(rseq.seen) >= rseq.maxIterations {
		return -1, false
	}

	val := rand.Intn(rseq.max)
	if val > rseq.min && !rseq.seen[val] {
		rseq.seen[val] = true
		return val, true
	}

	return rseq.next()
}

func NewPixelCursor(max, index int) *PixelCursor {
	return &PixelCursor{
		&SequentialSequence{
			index: index,
			max:   max,
		},
	}
}

func NewRandomizedPixelCursor(u UnderbyteImage, initPosition, maxIterations int) *PixelCursor {
	seed := u.randomizationSeed()
	rand.Seed(seed)

	return &PixelCursor{
		&RandomizedSequence{
			max:           u.pixelCount(),
			min:           initPosition,
			maxIterations: maxIterations,
			seen:          make(map[int]bool),
		},
	}
}
