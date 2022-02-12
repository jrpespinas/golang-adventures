package main

import "fmt"

func abs(x int) int {
	if x < 0 {
		return x * -1
	} else {
		return x
	}
}

func min(x, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}

func main() {
	var curr, next, min1, min2 int
	for true {
		fmt.Scanln(&curr, &next)
		if curr == -1 && next == -1 {
			break
		}
		curr++
		next++
		min1 = abs(curr - next)
		min2 = abs(100 - (curr + next))
		println(min(min1, min2))
	}
}
