package main

import "fmt"

func main() {
	count := 1
	for true {
		var chant string
		fmt.Scanln(&chant)
		if chant == "Hajj" {
			fmt.Printf("Case %v: %v\n", count, "Hajj-e-Akbar")
		} else if chant == "Umrah" {
			fmt.Printf("Case %v: %v\n", count, "Hajj-e-Asghar")
		} else if chant == "*" {
			break
		}
		count++
	}
}
