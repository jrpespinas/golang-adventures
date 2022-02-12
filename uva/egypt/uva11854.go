package main

import "fmt"

func main() {
	var a, b, c int
	for true {
		fmt.Scan(&a, &b, &c)
		if c == 0 && a == 0 && b == 0 {
			break
		}
		if c*c == a*a+b*b {
			println("right")
		} else {
			println("wrong")
		}
	}
}
