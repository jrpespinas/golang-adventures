package main

import "fmt"

func main() {
	var t, n, s, max int
	fmt.Scan(&t)
	for i := 1; i <= t; i++ {
		fmt.Scan(&n)
		max = 0
		for j := 0; j < n; j++ {
			fmt.Scan(&s)
			if s > max {
				max = s
			}
		}
		fmt.Printf("Case %v: %v\n", i, max)
	}
}
