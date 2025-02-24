package huffman

import "container/list"

// Algorithm from https://create.stephan-brumme.com/length-limited-prefix-codes/#package-merge

type node struct {
	symbol    int
	frequency uint64
}

func insertNode(lst *list.List, value *node) {
	for el := lst.Front(); el != nil; el = el.Next() {
		if el.Value.(*node).frequency >= value.frequency {
			lst.InsertBefore(value, el)
			return
		}
	}

	lst.PushBack(value)
}

func merge(lst *list.List, frequencies []uint64) {
	for symbol, frequency := range frequencies {
		if frequency > 0 {
			insertNode(lst, &node{
				symbol:    symbol + 1,
				frequency: frequency,
			})
		}
	}
}

func pack(node1, node2 *node) *node {
	return &node{
		symbol:    0, // No symbol for packed nodes
		frequency: node1.frequency + node2.frequency,
	}
}

// Computes the code lengths for symbols based on their frequencies.
func GenerateLengths(frequencies []uint64, maxLength uint8) []uint8 {
	// Find number of non zero frequencies
	nonZeroCount := 0
	for _, frequency := range frequencies {
		if frequency > 0 {
			nonZeroCount++
		}
	}

	lengths := make([]uint8, len(frequencies))

	// Edge case for empty codes
	if nonZeroCount == 0 {
		return lengths
	}

	// Edge case for 1 code
	if nonZeroCount == 1 {
		for symbol, frequency := range frequencies {
			if frequency > 0 {
				lengths[symbol] = 1
				return lengths
			}
		}
	}

	// Phase 1
	lists := make([]*list.List, maxLength)
	lists[0] = list.New()
	merge(lists[0], frequencies)

	for i := 1; i < int(maxLength); i++ {
		lists[i] = list.New()

		// Package
		el1 := lists[i-1].Front()
		for el1 != nil {
			el2 := el1.Next()
			if el2 == nil {
				break
			}

			lists[i].PushBack(pack(el1.Value.(*node), el2.Value.(*node)))
			el1 = el2.Next()
		}

		// Merge
		merge(lists[i], frequencies)
	}

	// Phase 2
	limit := (nonZeroCount - 1) * 2
	for i := len(lists) - 1; i >= 0 && limit != 0; i-- {
		packedCount := 0
		el := lists[i].Front()
		for ; limit > 0; limit-- {
			if el == nil {
				break
			}
			symbol := el.Value.(*node).symbol

			if symbol == 0 {
				packedCount++
			} else {
				lengths[symbol-1]++
			}

			el = el.Next()
		}

		limit = packedCount * 2
	}

	return lengths
}
