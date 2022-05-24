package router

import (
	"net/http"

	api "book-list/controller"
	"book-list/database"

	"github.com/go-chi/chi/v5"
)

func Router(database database.Storage) http.Handler {
	r := chi.NewRouter()

	// Routes for `books` resource
	r.Route("/api/v2/books", func(r chi.Router) {
		r.Post("/", api.CreateBookController(database))
		r.Get("/", api.GetBooksController(database))
		r.Get("/finished-books", api.GetFinishedBooksController(database))
		r.Get("/unfinished-books", api.GetUnfinishedBooksController(database))

		// Sub route for `books`
		// TODO: Add middleware to check for `id`
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", api.GetOneBookController(database))
			r.Put("/", api.EditBookController(database))
			r.Delete("/", api.DeleteBookController(database))
		})
	})

	// Routes `User` authentication
	r.Route("/api/v2", func(r chi.Router) {
		r.Post("/signup", api.Signup)
		r.Post("/login", api.Login)
		r.Get("/logout", api.Logout)
	})

	return r
}
