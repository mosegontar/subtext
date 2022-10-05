package underbyte

import (
	"reflect"
	"strings"
	"testing"
)

var testcases = []struct {
	message        []byte
	expectedOffset int
	expectedEnd    int
	expectedSize   int
	expectedBytes  []byte
}{
	{
		message:        []byte{},
		expectedOffset: int(2),
		expectedSize:   int(0),
		expectedBytes:  []byte{0, 0, 0, 0},
	},
	{
		message:        []byte("a"),
		expectedOffset: int(2),
		expectedSize:   int(1),
		expectedBytes:  []byte{0, 0, 0, 1},
	},
	{
		message:        []byte(strings.Repeat("a", 255)),
		expectedOffset: int(2),
		expectedSize:   int(255),
		expectedBytes:  []byte{0, 0, 0, 255},
	},
	{
		message:        []byte(strings.Repeat("a", 256)),
		expectedOffset: int(2),
		expectedSize:   int(256),
		expectedBytes:  []byte{0, 0, 1, 0},
	},
	{
		message:        []byte(strings.Repeat("😇", 65432)),
		expectedOffset: int(2),
		expectedSize:   int(261728),
		expectedBytes:  []byte{0, 3, 254, 96},
	},
}

func TestBytes(t *testing.T) {
	for _, testcase := range testcases {
		header := newHeader(testcase.message)

		expected := testcase.expectedBytes
		actual := header.Bytes()

		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
	}
}

func TestMessageOffset(t *testing.T) {
	for _, testcase := range testcases {
		header := newHeader(testcase.message)
		expected := testcase.expectedOffset
		actual := header.messageOffset()

		if expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
	}
}
