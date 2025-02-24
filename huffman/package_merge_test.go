package huffman

import (
	"reflect"
	"testing"
)

func TestGenerateLengths(t *testing.T) {
	tests := []struct {
		name        string
		frequencies []uint64
		maxLength   uint8
		expected    []uint8
	}{
		{"Empty frequencies", []uint64{}, 5, []uint8{}},
		{"Single symbol", []uint64{0, 10, 0}, 5, []uint8{0, 1, 0}},
		{"Two symbols with same frequency", []uint64{5, 5}, 5, []uint8{1, 1}},
		{"Regular case", []uint64{3, 1, 2}, 5, []uint8{1, 2, 2}},
		{"Max length constraint", []uint64{5, 3, 2, 1}, 2, []uint8{2, 2, 2, 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateLengths(tt.frequencies, tt.maxLength)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("For %s, expected %v but got %v", tt.name, tt.expected, result)
			}
		})
	}
}
