package dispatcher

type Process interface {
	CompareEntries(interface{}, interface{})
}

type Leaderboard struct {
	Entries []interface{}
	Rank    []int
}
