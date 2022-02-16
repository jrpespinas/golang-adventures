package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	mydb "server/database"
	"sync"
)

func main() {
	db := GetDatabase("data/data.json")
	http.ListenAndServe(":8080", db.Handler())
}

func GetDatabase(filename string) *mydb.Database {
	// Open json file
	log.Print("[GetDatabase] Opening JSON file")
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("[GetDatabase] unable to open the json file: %s", err.Error())
	}
	defer file.Close()

	// Read json file
	log.Print("[GetDatabase] Reading JSON file")
	byteData, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("[GetDatabase] unable to read json file: %s", err.Error())
	}

	// Parse json file
	log.Print("[GetDatabase] Parsing JSON file")
	var database *mydb.Database
	if err := json.Unmarshal([]byte(byteData), &database); err != nil {
		log.Fatalf("[GetDatabase] unable to unmarshal json file: %s", err.Error())
	}
	database.Mu = sync.Mutex{}

	log.Print("[GetDatabase] Database successfully started!")
	return database
}
