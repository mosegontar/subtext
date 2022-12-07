package underbyte

import (
	"reflect"
	"strings"
	"testing"
)

func Test_messageHeader(t *testing.T) {
	t.Run("buff contains expected values", func(t *testing.T) {
		testcases := []struct {
			message  []byte
			expected []byte
		}{
			{[]byte("hello"), []byte{0, 0, 0, 5}},
			{[]byte(strings.Repeat("a", 255)), []byte{0, 0, 0, 255}},
			{[]byte(strings.Repeat("a", 256)), []byte{0, 0, 1, 0}},
		}

		for _, testcase := range testcases {
			actual := messageHeader([]byte(testcase.message))
			expected := testcase.expected

			if !reflect.DeepEqual(actual, expected) {
				t.Errorf("got %v expected %v", actual, expected)
			}

		}
	})
}

func Test_headerBytesToInt(t *testing.T) {
	t.Run("returns the correct integer", func(t *testing.T) {
		testcases := []struct {
			headerBytes []byte
			val         int
		}{
			{[]byte{0, 0, 0, 5}, 5},
			{[]byte{0, 0, 0, 255}, 255},
			{[]byte{0, 0, 1, 0}, 256},
			{[]byte{0, 0, 1, 255}, 511},
			{[]byte{0, 0, 2, 0}, 512},
		}

		for _, testcase := range testcases {
			actual := headerBytesToInt([]byte(testcase.headerBytes))
			expected := testcase.val

			if !reflect.DeepEqual(actual, expected) {
				t.Errorf("got %v expected %v", actual, expected)
			}

		}

	})
}
