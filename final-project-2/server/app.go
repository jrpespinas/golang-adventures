package main

import (
	config "book-list/config"
	route "book-list/routes"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	// Get port number
	port := config.GetPort(os.Getenv("PORT"))

	log.Infof("SERVER: Listening at port%v", port)
	http.ListenAndServe(port, route.Router())
}
