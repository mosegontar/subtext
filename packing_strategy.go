package underbyte

type BytePacker struct {
	strategy PackingStrategy
}

type PackingStrategy interface {
	pack(*UnderbyteImage, []byte)
	messageTooLarge(*UnderbyteImage, []byte) bool
}

type DoublePackStrategy struct{}
type SinglePackStrategy struct{}
