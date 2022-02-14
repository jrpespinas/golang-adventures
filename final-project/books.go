package main

import (
	mydb "books/models"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

func main() {
	db := getDatabase("data/data.json")
	http.ListenAndServe(":8080", db.Handler())
}

func getDatabase(filename string) *mydb.Database {
	// Open json file
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("unable to open the json file: %s", err.Error())
	}
	defer file.Close()

	// Read json file
	byteData, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("unable to read json file: %s", err.Error())
	}

	// Parse json file
	var database *mydb.Database
	if err := json.Unmarshal([]byte(byteData), &database); err != nil {
		log.Fatalf("unable to unmarshal json file: %s", err.Error())
	}
	database.Mu = sync.Mutex{}
	return database
}
