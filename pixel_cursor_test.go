package underbyte

import (
	"reflect"
	"testing"
)

func Test_next(t *testing.T) {
	t.Run("it emits the expected values", func(t *testing.T) {
		cursor := NewPixelCursor(3, 0)
		type positionsResults struct {
			int
			bool
		}

		expected := []positionsResults{
			{0, true},
			{1, true},
			{2, true},
			{-1, false},
		}

		actual := []positionsResults{}
		for i := 0; i < 4; i++ {
			pos, ok := cursor.next()
			actual = append(actual, positionsResults{pos, ok})
		}

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("expected %v actual %v", expected, actual)
		}
	})
}
