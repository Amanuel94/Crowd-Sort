package selector

import (
	"network/shared"
	"network/shared/interfaces"
	"network/utils"

	"testing"

	"github.com/google/uuid"
)

func TestSelector(t *testing.T) {

	// sample input
	u := []interfaces.Comparable[int]{}
	size := utils.RandInt(2, 100)
	for i := 0; i < size; i++ {
		num := shared.NewInt(utils.RandInt(0, 100))
		u = append(u, shared.NewIndexedItem(num).(shared.IndexedItem[int]))
	}

	m := make(map[uuid.UUID]int)
	for ind, item := range u {
		m[item.GetIndex().(uuid.UUID)] = ind
	}

	// create selector
	cfg := NewConfig()
	s := NewSelector[int](*cfg)
	s.CreateGraph(u)
	p, ok := s.Next()

	for ok {
		i, j := p.F.GetIndex().(uuid.UUID), p.S.GetIndex().(uuid.UUID)
		pi := u[m[i]]
		pj := u[m[j]]

		if pi.Compare(pj) > 0 {
			u[m[i]], u[m[j]] = u[m[j]], u[m[i]]
		}
		s.PrepareNeighbours(p.Id)
		p, ok = s.Next()
	}

	// check if u is sorted
	for i := 0; i < size-1; i++ {
		if u[i].Compare(u[i+1]) > 0 {
			t.Fail()
		}
	}

}
func TestSelectorMultiple(t *testing.T) {
	n_test := 1000
	for i := 0; i < n_test; i++ {
		TestSelector(t)
	}
	t.Logf("Passed %d tests\n", n_test)
}
