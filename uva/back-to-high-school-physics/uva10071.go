package main

import (
	"fmt"
)

func main() {
	var v, t int

	_, err := fmt.Scan(&v, &t)
	for err == nil {
		fmt.Println((v * t) * 2)
		_, err = fmt.Scan(&v, &t)
	}
}
