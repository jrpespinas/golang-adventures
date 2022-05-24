package controllers

import (
	"encoding/json"
	"net/http"

	"book-list/database"
	"book-list/models"

	"github.com/go-chi/chi/v5"
)

// Return all book records
// Route: GET /api/v2/books
// Access: protected
func GetBooksController(db database.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		var Books []models.Book

		response := db.Database.GetAllBooks(Books)

		json.NewEncoder(w).Encode(response)
	}
}

// Create a book record
// Route: POST /api/v2/books
// Access: protected
func CreateBookController(db database.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")

		var Book models.Book
		json.NewDecoder(r.Body).Decode(&Book)

		response := db.Database.CreateBook(Book)

		json.NewEncoder(w).Encode(response)
	}
}

// Return book by id
// Route: GET /api/v2/books/{id}
// Access: protected
func GetOneBookController(db database.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		bookID := chi.URLParam(r, "id")

		response := db.Database.GetOneBook(bookID)

		json.NewEncoder(w).Encode(response)
	}
}

// Edit book by id
// Route: PUT /api/v2/books/{id}
// Access: protected
func EditBookController(db database.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		bookID := chi.URLParam(r, "id")

		var book models.Book
		json.NewDecoder(r.Body).Decode(&book)

		response := db.Database.EditBook(bookID, book)

		json.NewEncoder(w).Encode(response)
	}
}

// Delete book by id
// Route: DELETE /api/v2/books/{id}
// Access: protected
func DeleteBookController(db database.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		bookID := chi.URLParam(r, "id")

		response := db.Database.DeleteBook(bookID)

		json.NewEncoder(w).Encode(response)
	}
}

// Return all finished books
// Route: GET /api/v2/finished-books
// Access: protected
func GetFinishedBooksController(db database.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		var Books []models.Book

		response := db.Database.GetFinishedBooks(Books)

		json.NewEncoder(w).Encode(response)
	}
}

// Return all unfinished books
// Route: GET /api/v2/unfinished-books
// Access: protected
func GetUnfinishedBooksController(db database.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		var Books []models.Book

		response := db.Database.GetUnfinishedBooks(Books)

		json.NewEncoder(w).Encode(response)
	}
}
