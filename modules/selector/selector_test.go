package selector

import (
	"github.com/Amanuel94/crowdsort/interfaces"
	"github.com/Amanuel94/crowdsort/shared"
	"github.com/Amanuel94/crowdsort/utils"

	"testing"

	"github.com/google/uuid"
)

func TestSelector(t *testing.T) {

	// sample input
	u := []interfaces.Comparable[int]{}
	size := utils.RandInt(8, 8)
	for i := 0; i < size; i++ {
		num := shared.NewInt(utils.RandInt(0, 100))
		u = append(u, shared.NewIndexedItem[int](num))
	}

	m := make(map[uuid.UUID]shared.IndexedItem[int])
	for _, item := range u {
		m[item.GetIndex().(uuid.UUID)] = item.(shared.IndexedItem[int])
	}

	// create selector
	cfg := NewConfig()
	s := NewSelector[int](*cfg)
	s.CreateGraph(u)
	p, ok := s.Next()

	for ok {
		i, j := p.F, p.S
		pi := m[i]
		pj := m[j]

		pv := pi.GetValue()
		qv := pj.GetValue()

		if pi.Compare(pj) > 0 { // pi > pj
			pi.SetValue(qv)
			pj.SetValue(pv)
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
