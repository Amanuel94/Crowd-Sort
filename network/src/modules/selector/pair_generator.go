package selector

import (
	"network/utils"
)

// pair generator for batcher even-odd-merge sort (BEMS)
// TODO: Change to iterative version
func BEMS_pairs_generator(n int, delta int, offest int) []*[]int{

	argue(n >= 2, "n must be greater than 1")
	if (utils.Bit_count(n) != 1){
		n = utils.NextPower(n)
	}
	if n == 2{
		return []*[]int{{delta*0 + offest, delta*1 + offest}}
	}
	pairs :=[]*[]int{}
	pairs = append(pairs, BEMS_pairs_generator(n/2, delta, offest)...)
	pairs = append(pairs, BEMS_pairs_generator(n/2, delta, offest + delta*n/2)...)
	pairs = append(pairs, BEMS_merge(n/4, delta, offest)...)
	pairs = append(pairs, BEMS_merge(n/4, delta, offest + delta)...)

	for i := range(n-1){
		pair := &[]int{offest + 1 + i*2*delta, offest + (i+1)*delta}
		pairs = append(pairs, pair)
	}
	return pairs


}

func BEMS_merge(n int, d int, o int) []*[]int{
	if n == 1{
		return []*[]int{{o, o + d}}
	}
	pairs := []*[]int{}

	pairs = append(pairs, BEMS_merge(n/2, 2*d, o)...)
	pairs = append(pairs, BEMS_merge(n/2, 2*d, o + d)...)
	for i := range(n-1){
		pair := &[]int{o + 1 + i*2*d, o + (i+1)*d}
		pairs = append(pairs, pair)
	}
	return pairs
}

