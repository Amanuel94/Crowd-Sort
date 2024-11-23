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
	size := utils.RandInt(4, 4)
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
	p, ok := s.Batch()

	for s.q.size > 0 {
		if !ok {
			break
		}
		t.Logf("s.q.size = %d\n", s.q.size)
		i, j := p.f.GetIndex().(uuid.UUID), p.s.GetIndex().(uuid.UUID)
		pi := u[m[i]]
		pj := u[m[j]]

		if pi.Compare(pj) > 0 {
			u[m[i]], u[m[j]] = u[m[j]], u[m[i]]
		}
		m[i], m[j] = m[j], m[i]
		s.PrepareNeighbours(p.id)
		p, ok = s.Batch()
	}

	// check if u is sorted
	for i := 0; i < size-1; i++ {
		if u[i].Compare(u[i+1]) > 0 {
			// t.Logf("u[%d] = %d, u[%d] = %d\n", i, u[i].GetValue(), i+1, u[i+1].GetValue())
			t.Fail()
		}
	}

}
