package selector

// test for the pair generator

import (
	"testing"

	"github.com/amanuel94/crowdsort/utils"
)

func TestBEMS_pairs_generator(t *testing.T) {

	n_test := 2

	for i := 0; i < n_test; i++ {
		size := utils.RandInt(1, 1<<10)
		rand_list := make([]int, size)
		for j := 0; j < size; j++ {
			rand_list[j] = utils.RandInt(1, 1<<10)
		}
		pairs := BEMS_pairs_generator(size, 1, 0, nil)

		for _, pair := range pairs {
			i, j := (pair)[0], (pair)[1]
			if max(i, j) < size && rand_list[i] > rand_list[j] {
				rand_list[i], rand_list[j] = rand_list[j], rand_list[i]
			}
		}
		// check if rand_list is sorted
		for k := 0; k < size-1; k++ {
			if rand_list[k] > rand_list[k+1] {
				t.Fail()

			}

		}
	}
}
