package selector

// test for the pair generator

import (
	"network/utils"
	"testing"
)

func TestBEMS_pairs_generator(t *testing.T) {

	n_test := 10
	
	for i := 0; i < n_test; i++ {
		size := 1 << utils.RandInt(1, 7)
		rand_list := make([]int, size)
		for j := 0; j < size; j++ {
			rand_list[j] = utils.RandInt(1, 100)
		}
		pairs := BEMS_pairs_generator(size, 1, 0)
		
		for _, pair := range pairs {
			i, j := (*pair)[0], (*pair)[1]
			if rand_list[i] > rand_list[j] {
				rand_list[i], rand_list[j] = rand_list[j], rand_list[i]
			}
		// check if rand_list is sorted
		for k := 0; k < size-1; k++ {
			if rand_list[k] > rand_list[k+1] {
				t.Errorf("Test Failed: List is not sorted")
		}
	}
}
	}
}