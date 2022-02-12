package main

import (
	"fmt"
	"sort"
)

type array []int

func isArrayEven(arr []int) bool {
	// return true if array size is even, otherwise return false
	return len(arr)%2 == 0
}

func (arr array) median() int {
	var median int

	// if array size is even, assign the average of the two middle values to `median`
	// otherwise, assign the value which separates the array into two equal parts
	if isArrayEven(arr) {
		a := len(arr) / 2
		median = (arr[a-1] + arr[a]) / 2
	} else {
		median = arr[len(arr)/2]
	}
	return median
}

func main() {
	// Declare variables
	var in, median int
	var nums array

	// Ask for input
	_, err := fmt.Scan(&in)

	// Append values to array
	for err == nil {
		nums = append(nums, in)

		// Sort the array
		sort.Ints(nums)

		// Get median
		median = nums.median()

		// Display median
		println(median)

		// Ask for another input
		_, err = fmt.Scan(&in)
	}
}
