package selector

import (
	"github.com/Amanuel94/crowdsort/interfaces"
	"github.com/Amanuel94/crowdsort/shared"
	"github.com/Amanuel94/crowdsort/utils"

	"testing"
)

func TestSelector(t *testing.T) {
	// Generate sample input
	u := generateSampleInput()

	// Create a map for quick access
	m := createMap(u)

	// Create and configure selector
	cfg := NewConfig()
	selector := NewSelector[int](*cfg)
	selector.CreateGraph(u)

	// Process pairs
	processPairs(selector, m)

	// Verify if the slice is sorted
	if !isSorted(u) {
		t.Fail()
	}
}

func generateSampleInput() []interfaces.Comparable[int] {
	size := utils.RandInt(2, 100)
	u := make([]interfaces.Comparable[int], size)
	for i := 0; i < size; i++ {
		num := shared.NewInt(utils.RandInt(0, 100))
		u[i] = shared.NewWire[int](num)
	}
	return u
}

func createMap(u []interfaces.Comparable[int]) map[string]shared.Wire[int] {
	m := make(map[string]shared.Wire[int])
	for _, item := range u {
		m[item.GetIndex().(string)] = item.(shared.Wire[int])
	}
	return m
}

func processPairs(selector *Selector[int], m map[string]shared.Wire[int]) {
	p, ok := selector.Next()
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
		selector.PrepareNeighbours(p.Id)
		p, ok = selector.Next()
	}
}

func isSorted(u []interfaces.Comparable[int]) bool {
	for i := 0; i < len(u)-1; i++ {
		if u[i].Compare(u[i+1]) > 0 {
			return false
		}
	}
	return true
}
func TestSelectorMultiple(t *testing.T) {
	n_test := 1000
	for i := 0; i < n_test; i++ {
		TestSelector(t)
	}
	t.Logf("Passed %d tests\n", n_test)
}
