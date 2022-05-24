package main

import (
	config "book-list/config"
	"book-list/database"
	route "book-list/routes"
	_ "book-list/utils/log"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	log "github.com/sirupsen/logrus"
)

func main() {
	// Get port number
	port := config.GetPort(os.Getenv("PORT"))

	log.Info("Connecting to database...")

	mongodb_uri := fmt.Sprintf("%v:%v", os.Getenv("MONGODB_URI"), os.Getenv("DATABASE_PORT"))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI(mongodb_uri))

	db := database.NewMongoStorage("book-list-db", client)
	storage := database.Storage{
		Database: db,
	}

	http.ListenAndServe(port, route.Router(storage))
}
