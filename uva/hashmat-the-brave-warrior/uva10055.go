package main

import (
	"fmt"
	"io"
)

func main() {
	var h, e int
	_, err := fmt.Scan(&h, &e)
	for err == nil {
		if err == io.EOF {
			break
		}
		if e > h {
			fmt.Println(e - h)
		} else {
			fmt.Println(h - e)
		}
		_, err = fmt.Scan(&h, &e)
	}
}
