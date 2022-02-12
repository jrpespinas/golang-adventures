package main

import "fmt"

func main() {
	var set int
	fmt.Scan(&set)

	for i := 0; i < set; i++ {
		var a, b, c int
		var tri string
		fmt.Scan(&a, &b, &c)
		if a == b && b == c {
			tri = "Equilateral"
		} else if a == b || a == c || c == b {
			tri = "Isosceles"
		} else if a+b > c && a+c > b && c+b > a {
			tri = "Scalene"
		} else {
			tri = "Invalid"
		}
		fmt.Printf("Case %v: %v\n", i+1, tri)
	}
}
