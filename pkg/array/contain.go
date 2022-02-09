package array

// Contains checks for val in list.
func Contains(list []int64, val int64) bool {
	if list == nil {
		return false
	}
	for _, item := range list {
		if item == val {
			return true
		}
	}
	return false
}
