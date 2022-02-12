package main

import "fmt"

func main() {
	var set int
	fmt.Scan(&set)

	for i := 0; i < set; i++ {
		var a, b int
		var sum int = 0

		fmt.Scan(&a, &b)

		for j := a; j < b; j++ {
			if j%2 != 0 {
				sum += j
			}
		}
		fmt.Println(sum)
	}
}
