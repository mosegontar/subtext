package underbyte

type PackingStrategy interface {
	pack(*UnderbyteImage, []byte, int) int
	unpack(*UnderbyteImage, int, int) []byte
}

type DoublePackStrategy struct{}
type SinglePackStrategy struct{}
