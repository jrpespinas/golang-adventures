package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Contacts struct {
	Last    string `json:"last"`
	First   string `json:"first"`
	Company string `json:"company"`
	Address string `json:"address"`
	Country string `json:"country"`
	Positon string `json:"position"`
}

type Database struct {
	ID        int
	mu        sync.Mutex
	IDPointer map[int]interface{}
}

func (db *Database) handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var id int
		if r.URL.Path == "/contacts" {
			db.processCollection(w, r)
		} else if n, _ := fmt.Sscanf(r.URL.Path, "/contacts/%d", &id); n == 1 {
			db.processID(id, w, r)
		} else {
			http.Error(w, "url accessed not found", http.StatusNotFound)
		}
	}
}

func (db *Database) processCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		// Decode json body
		var item Contacts
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		// Add record to database
		db.mu.Lock()
		db.IDPointer[db.ID] = item
		db.ID++
		db.mu.Unlock()

		// Set http status code
		w.WriteHeader(http.StatusCreated)

	case "GET":
		if len(db.IDPointer) == 0 {
			w.WriteHeader(http.StatusNoContent)
		}
		// Get all records
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(db.IDPointer); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "PUT":
		w.WriteHeader(http.StatusMethodNotAllowed)
	case "DELETE":
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (db *Database) processID(id int, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		w.WriteHeader(http.StatusMethodNotAllowed)
	case "GET":
		if _, ok := db.IDPointer[id]; ok {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(db.IDPointer[id]); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "contact ID does not exist", http.StatusNotFound)
		}
	case "PUT":
		if _, ok := db.IDPointer[id]; ok {
			// Decode json body
			var item Contacts
			if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			// Update record to database
			db.mu.Lock()
			db.IDPointer[id] = item
			db.mu.Unlock()

			// Set http status code
			w.WriteHeader(http.StatusOK)
		}
	case "DELETE":
		db.mu.Lock()
		if _, ok := db.IDPointer[id]; ok {
			delete(db.IDPointer, id)
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "contact ID does not exist", http.StatusNotFound)
		}
		db.mu.Unlock()
	}
}
func main() {
	db := &Database{IDPointer: make(map[int]interface{})}
	http.ListenAndServe(":8080", db.handler())
}
