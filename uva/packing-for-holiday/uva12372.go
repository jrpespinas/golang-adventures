package main

import "fmt"

func main() {
	var n, l, w, h int
	fmt.Scan(&n)
	for i := 1; i <= n; i++ {
		fmt.Scan(&l, &w, &h)
		if l > 20 || w > 20 || h > 20 {
			fmt.Printf("Case %v: %v\n", i, "bad")
		} else {
			fmt.Printf("Case %v: %v\n", i, "good")
		}
	}
}
