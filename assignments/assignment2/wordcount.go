package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
)

type SafeCounter struct {
	mu sync.Mutex
	v  map[string]int
}

func (c *SafeCounter) Increment(key string, value int) {
	c.mu.Lock()
	c.v[key] += value
	c.mu.Unlock()
}

func (c *SafeCounter) Value(key string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.v[key]
}

func main() {
	// Get files
	files := os.Args[1:]
	numOfFiles := len(files)
	if numOfFiles == 0 {
		log.Fatalf("No files selected")
	}

	// set regular expression
	re, err := regexp.Compile(`[^\w]`)
	if err != nil {
		log.Fatalf("this regular expression is invalid: %s", err)
	}

	// set counter
	counter := SafeCounter{v: make(map[string]int)}
	// counterPtr := &counter

	// initialize array
	var wordsArray []string

	// initialize channel
	channel := make(chan map[string]int)

	// Store words to array per file
	i := 0
	for i < numOfFiles {
		arr, err := getWordsPerFile(files[i], re)
		if err != nil {
			log.Fatalf("scanner error: %s", err)
		}
		wordsArray = append(wordsArray, arr...)
		i++
	}

	// concurrently count words
	length := len(wordsArray)
	go wordCount(wordsArray[0:(length/2)], channel)
	go wordCount(wordsArray[(length/2):], channel)

	for i := 0; i < 2; i++ {
		words := <-channel
		for k, v := range words {
			counter.Increment(k, v)
		}
	}

	printMap(counter.v)
}

func wordCount(arr []string, ch chan map[string]int) {
	frequency := map[string]int{}
	for _, word := range arr {
		frequency[word]++
	}
	ch <- frequency
}

func getWordsPerFile(filename string, re *regexp.Regexp) ([]string, error) {
	// Open file
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	// Set Scanner
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	// initialize local array
	var arr []string

	// Words to hash map
	for scanner.Scan() {
		word := (normalizeWord(scanner.Text(), re))
		arr = append(arr, word)
	}

	return arr, scanner.Err()
}

// A function to clean the string before adding to the hash map
func normalizeWord(word string, re *regexp.Regexp) string {
	lowered := strings.ToLower(word)
	normalizedWord := re.ReplaceAllString(lowered, "")
	return normalizedWord
}

// Print contents of hashmap alphabetically
func printMap(wordcount map[string]int) {
	keys := make([]string, 0, len(wordcount))

	for k := range wordcount {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Println(k, wordcount[k])
	}
}
