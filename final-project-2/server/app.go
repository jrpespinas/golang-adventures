package main

import (
	config "book-list/config"
	route "book-list/routes"
	"log"
	"net/http"
	"os"
)

func main() {
	// Get port number
	port := config.GetPort(config.Getenv("PORT"))

	// Start Logger
	f, err := os.OpenFile("logs/logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("ERROR: %v", err.Error())
	}
	log.SetOutput(f)
	defer f.Close()

	log.Printf("SERVER: Listening at port%v", port)
	http.ListenAndServe(port, route.Router())
}
