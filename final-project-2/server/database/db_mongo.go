package database

import (
	"book-list/models"
	"context"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoStorage struct {
	dbname string
	conn   *mongo.Client
}

func NewMongoStorage(dbname string, conn *mongo.Client) *MongoStorage {
	return &MongoStorage{
		dbname: dbname,
		conn:   conn,
	}
}

func (mongodb MongoStorage) CreateBook(book models.Book) models.Response {
	collection := mongodb.conn.Database(mongodb.dbname).Collection("book")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, err := collection.InsertOne(ctx, book)
	if err != nil {
		return models.Response{
			Status:  "Fail",
			Message: err,
		}
	} else {
		return models.Response{
			Status:  "Success",
			Message: result,
		}
	}
}

func (mongodb MongoStorage) GetOneBook(id string) models.Response {
	collection := mongodb.conn.Database(mongodb.dbname).Collection("book")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	objID, _ := primitive.ObjectIDFromHex(id)
	var book models.Book
	err := collection.FindOne(ctx, models.Book{ID: objID}).Decode(&book)
	if err != nil {
		return models.Response{
			Status: "Fail",
			Message: models.ErrorMessage{
				Code:  http.StatusNotFound,
				Error: err.Error(),
			},
		}
	} else {
		return models.Response{
			Status:  "Success",
			Message: book,
		}
	}
}

func (mongodb MongoStorage) GetAllBooks(books []models.Book) models.Response {
	collection := mongodb.conn.Database(mongodb.dbname).Collection("book")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return models.Response{
			Status:  "Fail",
			Message: err,
		}
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var book models.Book
		if err = cursor.Decode(&book); err != nil {
			return models.Response{
				Status:  "Fail",
				Message: err,
			}
		}
		books = append(books, book)
	}

	return models.Response{
		Status:  "Success",
		Message: books,
	}
}

func (mongodb MongoStorage) GetFinishedBooks(books []models.Book) models.Response {
	collection := mongodb.conn.Database(mongodb.dbname).Collection("book")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{"is_finished": true})
	if err != nil {
		return models.Response{
			Status:  "Fail",
			Message: err,
		}
	}
	if err = cursor.All(ctx, &books); err != nil {
		return models.Response{
			Status:  "Fail",
			Message: err,
		}
	}

	return models.Response{
		Status:  "Success",
		Message: books,
	}
}

func (mongodb MongoStorage) GetUnfinishedBooks(books []models.Book) models.Response {
	collection := mongodb.conn.Database(mongodb.dbname).Collection("book")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{"is_finished": false})
	if err != nil {
		return models.Response{
			Status:  "Fail",
			Message: err,
		}
	}
	if err = cursor.All(ctx, &books); err != nil {
		return models.Response{
			Status:  "Fail",
			Message: err,
		}
	}

	return models.Response{
		Status:  "Success",
		Message: books,
	}
}

func (mongodb MongoStorage) DeleteBook(id string) models.Response {
	collection := mongodb.conn.Database(mongodb.dbname).Collection("book")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	objID, _ := primitive.ObjectIDFromHex(id)

	filter := bson.M{"_id": objID}
	result, _ := collection.DeleteOne(ctx, filter)

	if result.DeletedCount == 0 {
		return models.Response{
			Status: "Fail",
			Message: models.ErrorMessage{
				Code:  http.StatusNotFound,
				Error: "Document not found",
			},
		}
	} else {
		return models.Response{
			Status:  "Success",
			Message: fmt.Sprintf("Book ID %v deleted", id),
		}
	}

}

func (mongodb MongoStorage) EditBook(id string, book models.Book) models.Response {
	collection := mongodb.conn.Database(mongodb.dbname).Collection("book")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	objID, _ := primitive.ObjectIDFromHex(id)

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": book}

	result, _ := collection.UpdateOne(ctx, filter, update)
	if result.MatchedCount == 0 {
		return models.Response{
			Status: "Fail",
			Message: models.ErrorMessage{
				Code:  http.StatusNotFound,
				Error: "Document not found",
			},
		}
	} else {
		return models.Response{
			Status:  "Success",
			Message: fmt.Sprintf("Book ID %v updated", id),
		}
	}

}
