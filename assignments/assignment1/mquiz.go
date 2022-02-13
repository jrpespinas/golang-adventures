package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
)

func check_file(e error) {
	if e != nil {
		log.Fatalf("failed to open csv file: %s", e)
	}
}

func check_content(e error) {
	if e != nil {
		log.Fatalf("failed to parse csv file: %s", e)
	}
}

func check_int_parser(e error) {
	if e != nil {
		panic(e)
	}
}

// Retrieves an array of random numbers (without repetition) which is drawn
// from the number of rows `length` of the csv. However, if `count` is larger
// than `length`, then we allow minimally the repetition of rows.
func get_random_rows(count *int, lines [][]string) []int {
	length := len(lines)
	var array, items, rand_array []int
	var running_count int = *count

	if length < *count {
		for running_count > 0 {
			if length > running_count {
				rand_array = rand.Perm(length)[:running_count]
				items = append(items, rand_array...)
				running_count = 0
				break
			}
			rand_array = rand.Perm(length)[:length]
			items = append(items, rand_array...)
			running_count = running_count - length
		}
	} else if *count <= 0 {
		log.Fatalf("choose 'n' greater than 0")
	} else {
		items = rand.Perm(length)[:*count]
	}

	array = append(array, items...)
	return array
}

func main() {
	// Get and parse flags
	filename := flag.String("csv", "problem.csv", "csv in format 'question,answer'")
	count := flag.Int("n", 10, "number of problems to solve")
	flag.Parse()

	// Open csv file
	f, err := os.Open(*filename)
	check_file(err)
	defer f.Close()

	// Get contents of csv file
	reader := csv.NewReader(f)
	lines, err := reader.ReadAll()
	check_content(err)

	// Initialize variables
	var score int = 0
	var answer int

	// Get random rows from the csv file
	items := get_random_rows(count, lines)

	// Choose random equations from csv
	for i := 0; i < *count; i++ {
		// User prompt
		fmt.Printf("%v = ", lines[items[i]][0])
		fmt.Scanf("%v", &answer)

		// parse answer
		result, err := strconv.Atoi(lines[items[i]][1])
		check_int_parser(err)

		// increment score if correct answer
		if answer == result {
			score++
		}
	}

	// Results
	fmt.Printf("You answered %v out of %v questions correclty", score, *count)
}
