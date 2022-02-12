package main

import (
	"fmt"
)

func main() {
	// initialize variables
	var v, t int

	_, err := fmt.Scan(&v, &t)
	// Ask for inputs
	for err == nil {
		fmt.Println((v * t) * 2)
		_, err = fmt.Scan(&v, &t)
	}
}
