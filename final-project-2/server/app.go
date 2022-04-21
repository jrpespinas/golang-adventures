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
	port := config.GetPort(os.Getenv("PORT"))

	log.Printf("SERVER: Listening at port%v", port)
	http.ListenAndServe(port, route.Router())
}
