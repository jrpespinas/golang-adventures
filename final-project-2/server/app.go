package main

import (
	config "book-list/config"
	route "book-list/routes"
	logs "book-list/utils/logs"
	"net/http"
	"os"
)

func main() {
	// Get port number
	port := config.GetPort(os.Getenv("PORT"))

	// Initialize 3rd party log
	logs.Log.Sugar().Infof("SERVER: Listening at port%v", port)
	http.ListenAndServe(port, route.Router())
}
