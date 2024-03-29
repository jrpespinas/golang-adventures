package controllers

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Return all book records
// Route: GET /api/v2/books
// Access: protected
func GetBooks(w http.ResponseWriter, r *http.Request) {
	log.Info("Getting all books")
	w.Write([]byte("Hello"))
}

// Create a book record
// Route: POST /api/v2/books
// Access: protected
func CreateBook(w http.ResponseWriter, r *http.Request) {
	log.Info("Creating a book")
	w.Write([]byte("Hello"))
}

// Return book by id
// Route: GET /api/v2/books/{id}
// Access: protected
func GetBook(w http.ResponseWriter, r *http.Request) {
	log.Info("Getting a book")
	w.Write([]byte("Hello"))
}

// Edit book by id
// Route: PUT /api/v2/books/{id}
// Access: protected
func EditBook(w http.ResponseWriter, r *http.Request) {
	log.Info("Editing a book")
	w.Write([]byte("Hello"))
}

// Delete book by id
// Route: DELETE /api/v2/books/{id}
// Access: protected
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	log.Info("Deleting a book")
	w.Write([]byte("Hello"))
}

// Return all finished books
// Route: GET /api/v2/finished-books
// Access: protected
func GetFinishedBooks(w http.ResponseWriter, r *http.Request) {
	log.Info("Getting finished books")
	w.Write([]byte("Hello"))
}

// Return all unfinished books
// Route: GET /api/v2/unfinished-books
// Access: protected
func GetUnfinishedBooks(w http.ResponseWriter, r *http.Request) {
	log.Info("Getting unfinished books")
	w.Write([]byte("Hello"))
}
