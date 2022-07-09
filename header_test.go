package underbyte

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

var testcases = []struct {
	message        []byte
	expectedOffset uint
	expectedEnd    uint
	expectedSize   uint
	expectedBytes  []byte
}{
	{
		message:        []byte{},
		expectedOffset: uint(1),
		expectedEnd:    uint(0),
		expectedSize:   uint(0),
		expectedBytes:  []byte{0},
	},
	{
		message:        []byte("a"),
		expectedOffset: uint(2),
		expectedEnd:    uint(2),
		expectedSize:   uint(1),
		expectedBytes:  []byte{1, 1},
	},
	{
		message:        []byte(strings.Repeat("a", 255)),
		expectedOffset: uint(2),
		expectedEnd:    uint(256),
		expectedSize:   uint(255),
		expectedBytes:  []byte{1, 255},
	},
	{
		message:        []byte(strings.Repeat("a", 256)),
		expectedOffset: uint(3),
		expectedEnd:    uint(258),
		expectedSize:   uint(256),
		expectedBytes:  []byte{2, 0, 1},
	},
	{
		message:        []byte(strings.Repeat("😇", 65432)),
		expectedOffset: uint(4),
		expectedEnd:    uint(261731),
		expectedSize:   uint(261728),
		expectedBytes:  []byte{3, 96, 254, 3},
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

func TestMessageSize(t *testing.T) {
	for _, testcase := range testcases {
		header := newHeader(testcase.message)
		expected := testcase.expectedSize
		actual := header.messageSize()

		if expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
	}
}

func TestMessageEnd(t *testing.T) {
	for _, testcase := range testcases {
		header := newHeader(testcase.message)
		expected := testcase.expectedEnd
		actual := header.messageEnd()

		if expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
	}
}

func TestBytesToInt(t *testing.T) {
	var testcases = []struct {
		input    []byte
		expected uint
	}{
		{[]byte{128}, 128},
		{[]byte{255}, 255},
		{[]byte{0, 1}, 256},
		{[]byte{1, 1}, 257},
		{[]byte{255, 1}, 511},
		{[]byte{0, 2, 0}, 512},
		{[]byte{8, 2, 0}, 520},
		{[]byte{1, 2, 1}, 66049},
		{[]byte{0, 0, 0}, 0},
		{[]byte{0, 0, 1}, 65536},
		{[]byte{19, 0, 1}, 65555},
		{[]byte{}, 0},
	}

	for _, testcase := range testcases {
		t.Run(fmt.Sprintf("%v", testcase.input), func(t *testing.T) {
			actual := bytesToInt(testcase.input)

			if testcase.expected != actual {
				t.Errorf("expected %d, got %d", testcase.expected, actual)
			}
		})
	}
}
