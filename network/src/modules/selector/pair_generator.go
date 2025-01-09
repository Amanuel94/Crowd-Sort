package selector

import (
	"network/utils"
)

// pair generator for batcher even-odd-merge sort (BEMS)
// TODO: Change to iterative version
func BEMS_pairs_generator(n int, d int, o int, msg *chan interface{}) [][]int {

	deferPanic(msg)
	argue(n >= 2, "n must be greater than 1")
	if utils.Bit_count(n) > 1 {
		n = utils.NextPower(n)
	}
	if n == 2 {
		return [][]int{{d*0 + o, d*1 + o}}
	}
	pairs := [][]int{}
	pairs = append(pairs, BEMS_pairs_generator(n/2, d, o, msg)...)
	pairs = append(pairs, BEMS_pairs_generator(n/2, d, o+d*n/2, msg)...)
	pairs = append(pairs, BEMS_merge(n/4, 2*d, o)...)
	pairs = append(pairs, BEMS_merge(n/4, 2*d, o+d)...)

	for i := range n/2 - 1 {
		pair := []int{o + d + (i)*2*d, o + (i+1)*2*d}
		pairs = append(pairs, pair)
	}
	return pairs

}

func BEMS_merge(n int, d int, o int) [][]int {
	if n == 1 {
		return [][]int{{o, o + d}}
	}
	pairs := [][]int{}

	pairs = append(pairs, BEMS_merge(n/2, 2*d, o)...)
	pairs = append(pairs, BEMS_merge(n/2, 2*d, o+d)...)
	for i := range n - 1 {
		pair := []int{o + d + i*2*d, o + (i+1)*2*d}
		pairs = append(pairs, pair)
	}
	return pairs
}
