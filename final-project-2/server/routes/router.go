package router

import (
	"net/http"

	api "book-list/controller"

	"github.com/go-chi/chi/v5"
)

func Router() http.Handler {
	r := chi.NewRouter()

	// Routes for `books` resource
	r.Route("/api/v2/books", func(r chi.Router) {
		r.Post("/", api.CreateBook)
		r.Get("/", api.GetBooks)
		r.Get("/finished-books", api.GetFinishedBooks)
		r.Get("/unfinished-books", api.GetUnfinishedBooks)

		// Sub route for `books`
		// TODO: Add middleware to check for `id`
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", api.GetBook)
			r.Put("/", api.EditBook)
			r.Delete("/", api.DeleteBook)
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
