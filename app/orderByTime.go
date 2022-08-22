package app

import "sort"

type orderByTime struct {
	// seconds
	t int64

	// index of original slice
	p int
}

// sortAndSet will sort the items by time in descending order.
// But, most of the elements in items are already sorted in ascending order.
// In order to increase the sorting, no need to reverse the items.
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
