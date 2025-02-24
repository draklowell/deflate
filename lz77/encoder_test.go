package lz77

import "testing"

func TestFindMatch(t *testing.T) {
	tests := []struct {
		name           string
		bufferIn       []byte
		bufferOut      []byte
		expectedDist   uint16
		expectedLength uint16
	}{
		{"Match at end", []byte("ABC"), []byte("CCABC"), 3, 3},
		{"Repeating pattern", []byte("ABACABACA"), []byte("ABAC"), 4, 9},
		{"No match", []byte("ABC"), []byte("FF"), 0, 0},
		{"Full repeat", []byte("XYZ"), []byte("XYZXYZ"), 6, 3},
		{"Partial repeat", []byte("XYZ"), []byte("XYZZXY"), 6, 3},
		{"Overlapping sequence", []byte("BAA"), []byte("AAAABAAA"), 4, 3},
		{"Long match", []byte("LOHEL"), []byte("HELLOHELLO"), 7, 5},
		{"Single character match", []byte("A"), []byte("A"), 0, 0},
		{"Single character no match", []byte("A"), []byte("B"), 0, 0},
		{"Empty input", []byte(""), []byte("XYZ"), 0, 0},
		{"Empty output", []byte("XYZ"), []byte(""), 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			length, distance := findMatch(
				NewTestRingBuffer(uint16(len(tt.bufferIn)), tt.bufferIn),
				NewTestRingBuffer(uint16(len(tt.bufferOut)), tt.bufferOut),
			)

			if length != tt.expectedLength || distance != tt.expectedDist {
				t.Errorf("%s: expected (length: %d, distance: %d), got (length: %d, distance: %d)",
					tt.name, tt.expectedLength, tt.expectedDist, length, distance)
			}
		})
	}
}
