package main

import (
	"iter"
	"time"

	"github.com/Amanuel94/crowdsort/interfaces"
	"github.com/Amanuel94/crowdsort/modules/io"
	"github.com/Amanuel94/crowdsort/shared"
	"github.com/Amanuel94/crowdsort/utils"
)

// Example Usage of the crowdsort package
func main() {

	n_items := 10
	n_cmps := 3
	items := generateItems(n_items)
	comparators := generateComparators(n_cmps)
	io_cfg := io.NewConfig(items, comparators, 2).WithBufferSize(10).WithTaskLimit(11)
	io := io.New(io_cfg)
	go io.StartDispatcher()
	go io.ShowLeaderboard()
	io.Wait()

}

func generateItems(n int) iter.Seq[*interfaces.Comparable[int]] {

	randIntArr := make([]int, n)
	for i := 0; i < n; i++ {
		randIntArr[i] = utils.RandInt(1, 100)
	}

	items := utils.SliceToSeq(randIntArr)
	return utils.Map(func(v int) *interfaces.Comparable[int] {
		item := shared.NewInt(v)
		return &item
	}, items)

}

func generateComparators(n int) iter.Seq[shared.CmpFunc[int]] {

	comparators := make([]shared.CmpFunc[int], n)

	for i := 0; i < n; i++ {
		index := i
		comparators[index] = func(a *interfaces.Comparable[int], b *interfaces.Comparable[int]) (int, error) {

			w := utils.RandInt(1, 3)
			waitTime := time.Duration(w) * time.Second
			time.Sleep(waitTime)
			return (*a).Compare(*b), nil
		}
	}
	return utils.SliceToSeq(comparators)

}
