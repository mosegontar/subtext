package underbyte

type PackingStrategy interface {
	pack(*UnderbyteImage, []byte, int) int
	messageTooLarge(*UnderbyteImage, []byte, []byte) bool
}

type DoublePackStrategy struct{}
type SinglePackStrategy struct{}
