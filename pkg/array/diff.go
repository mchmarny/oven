package array

// GetDiff returns items from b that are NOT in a.
func GetDiff(a, b []int64) (diff []int64) {
	m := make(map[int64]bool)

	for _, item := range a {
		m[item] = true
	}

	for _, item := range b {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}
