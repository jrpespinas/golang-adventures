package main

import (
	config "book-list/config"
	route "book-list/routes"
	_ "book-list/utils/log"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	log "github.com/sirupsen/logrus"
)

func main() {
	// Get port number
	port := config.GetPort(os.Getenv("PORT"))

	log.Info("Connecting to database...")

	mongodb_uri := fmt.Sprintf("%v:%v", os.Getenv("MONGODB_URI"), os.Getenv("DATABASE_PORT"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodb_uri))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Panic(err)
		}
	}()

	/*
	   List databases
	*/
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Database Found")
	log.Info(databases)

	log.Infof("Listening at port%v", port)
	http.ListenAndServe(port, route.Router())
}
