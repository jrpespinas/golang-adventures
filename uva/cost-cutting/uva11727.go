package main

import (
	"fmt"
	"sort"
)

func main() {
	var set int
	fmt.Scan(&set)

	for i := 0; i < set; i++ {
		var a, b, c int
		fmt.Scan(&a, &b, &c)
		salary := []int{a, b, c}
		sort.Ints(salary)
		fmt.Printf("Case %v: %v\n", i+1, salary[1])
	}
}
