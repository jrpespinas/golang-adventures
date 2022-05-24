package database

import (
	"book-list/models"
)

type DB interface {
	CreateBook(book models.Book) models.Response
	GetAllBooks(books []models.Book) models.Response
	GetOneBook(id string) models.Response
	GetFinishedBooks(books []models.Book) models.Response
	GetUnfinishedBooks(books []models.Book) models.Response
	DeleteBook(id string) models.Response
	EditBook(id string, book models.Book) models.Response
}

type Storage struct {
	Database DB
}
