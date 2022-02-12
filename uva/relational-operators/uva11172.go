package main

import "fmt"

func main() {
	var set int
	fmt.Scan(&set)

	for i := 0; i < set; i++ {
		var a, b int
		fmt.Scan(&a, &b)
		switch {
		case a < b:
			fmt.Println("<")
		case a > b:
			fmt.Println(">")
		default:
			fmt.Println("=")
		}
	}
}
