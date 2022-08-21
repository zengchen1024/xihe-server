package app

import "sort"

type orderByTime struct {
	// seconds
	t int64

	// index of original slice
	p int
}

// Most of the elements in items are already sorted in ascending order.
// No need to reverse the items in order to increase the sorting.
func sortAndSet(items []orderByTime, do func(int, int) error) error {
	sort.Slice(items, func(i, j int) bool {
		return items[i].t < items[j].t
	})

	for i, j := len(items)-1, 0; i >= 0; i-- {
		if err := do(items[i].p, j); err != nil {
			return err
		}

		j++
	}

	return nil
}
