package huffman

import (
	"reflect"
	"testing"
)

func TestGenerateCodes(t *testing.T) {
	tests := []struct {
		name     string
		lengths  []uint8
		expected []uint32
	}{
		{"Empty input", []uint8{}, []uint32{}},
		{"Single symbol", []uint8{3}, []uint32{0}},
		{"Two symbols with different lengths", []uint8{1, 2}, []uint32{0, 2}},
		{"Multiple symbols with same lengths", []uint8{2, 2, 2}, []uint32{0, 1, 2}},
		{"Regular case", []uint8{3, 3, 3, 3, 2, 4}, []uint32{2, 3, 4, 5, 0, 12}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			codes := generateCodes(tt.lengths)
			if !reflect.DeepEqual(codes, tt.expected) {
				t.Errorf("generateCodes(%v) = %v, want %v", tt.lengths, codes, tt.expected)
			}
		})
	}
}
