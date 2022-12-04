package underbyte

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
	index int
	max   int
	count int

	seed int64
	seen map[int]bool
}

type PixelCursor struct {
	Sequence
}

//func (pc *PixelCursor) next() (int, bool) {
//	if pc.position > pc.max-1 {
//		return -1, false
//	}
//
//	current := pc.position
//	pc.position++
//	return current, true
//}
//

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

//func (rseq RandomizedSequence) next() (int, bool) {
//	if len(rseq.seen) >= rseq.max {
//		panic("help")
//	}
//	return
//}

func NewPixelCursor(max int, index int) *PixelCursor {
	return &PixelCursor{
		&SequentialSequence{
			index: index,
			max:   max,
		},
	}
}
