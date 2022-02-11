package cards

// GetItem retrieves an item from a slice at given position. The second return value indicates whether
// the given index exists in the slice or not.
func GetItem(slice []int, index int) (int, bool) {
	if index < len(slice) && index > -1 {
		return slice[index], true
	} else {
		return 0, false
	}
}

// SetItem writes an item to a slice at given position overwriting an existing value.
// If the index is out of range the value needs to be appended.
func SetItem(slice []int, index, value int) []int {
	if index < len(slice) && index > -1 {
		slice[index] = value
		return slice
	} else {
		slice = append(slice, value)
		return slice
	}
}

// PrefilledSlice creates a slice of given length and prefills it with the given value.
func PrefilledSlice(value, length int) []int {
	if length > -1 {
		arr := make([]int, length)
		for i := 0; i < length; i++ {
			arr[i] = value
		}
		return arr
	} else {
		return []int{}
	}
}

// RemoveItem removes an item from a slice by modifying the existing slice.
func RemoveItem(slice []int, index int) []int {
	if _, ok := GetItem(slice, index); ok {
		return append(slice[:index], slice[index+1:]...)
	}
	return slice
}
