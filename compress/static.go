package compress

import "deflate/huffman"

const (
	lengthMax   = 258
	distanceMax = 32768
)

var (
	lengthOffsets   = []uint16{3, 4, 5, 6, 7, 8, 9, 10, 11, 13, 15, 17, 19, 23, 27, 31, 35, 43, 51, 59, 67, 83, 99, 115, 131, 163, 195, 227, 258}
	lengthExtraBits = []uint8{0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 0}

	distanceOffsets   = []uint16{1, 2, 3, 4, 5, 7, 9, 13, 17, 25, 33, 49, 65, 97, 129, 193, 257, 385, 513, 769, 1025, 1537, 2049, 3073, 4097, 6145, 8193, 12289, 16385, 24577}
	distanceExtraBits = []uint8{0, 0, 0, 0, 1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 6, 7, 7, 8, 8, 9, 9, 10, 10, 11, 11, 12, 12, 13, 13}
)

func getStaticLetterLengths() [288]uint8 {
	letterLengths := [288]uint8{}
	for i := range 288 {
		if i <= 143 {
			letterLengths[i] = 8
		} else if i <= 255 {
			letterLengths[i] = 9
		} else if i <= 279 {
			letterLengths[i] = 7
		} else {
			letterLengths[i] = 8
		}
	}

	return letterLengths
}

func getStaticDistanceLengths() [32]uint8 {
	distanceLengths := [32]uint8{}
	for i := range 32 {
		distanceLengths[i] = 5
	}

	return distanceLengths
}

var distanceLengths = getStaticDistanceLengths()
var letterLengths = getStaticLetterLengths()
var (
	staticLetterEncodingTree, _   = huffman.NewEncodingTree(letterLengths[:])
	staticDistanceEncodingTree, _ = huffman.NewEncodingTree(distanceLengths[:])
	staticLetterDecodingTree, _   = huffman.NewDecodingTree(letterLengths[:])
	staticDistanceDecodingTree, _ = huffman.NewDecodingTree(distanceLengths[:])
)
