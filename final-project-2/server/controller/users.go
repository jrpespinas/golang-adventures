package controllers

import (
	"log"
	"net/http"
)

// Create a user account
// Route: POST /api/v2/signup
// Access: public
func Signup(w http.ResponseWriter, r *http.Request) {
	log.Print("[Sign up]")
	w.Write([]byte("Hello"))
}

// Log in with email and password
// Route: POST /api/v2/login
// Access: public
func Login(w http.ResponseWriter, r *http.Request) {
	log.Print("[Log in]")
	w.Write([]byte("Hello"))
}

// Log out from the app
// Route: GET /api/v2/logout
// Access: public
func Logout(w http.ResponseWriter, r *http.Request) {
	log.Print("[Log out]")
	w.Write([]byte("Hello"))
}
