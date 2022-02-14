package models

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("secret_key")

func (db *Database) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// var bookID int
		if r.URL.Path == "/login" {
			db.Login(w, r)
		} else if r.URL.Path == "/signup" {
			db.SignUp(w, r)
		} else if r.URL.Path == "/logout" {
			db.LogOut(w, r)
		} else {
			http.Error(w, "url does not exist", http.StatusNotImplemented)
		}

	}
}

///////////////////////////////////////////////////////////////
// LOGGING FEATURES
///////////////////////////////////////////////////////////////

func (db *Database) Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		log.Print("Creating variables")
		var user User

		// Get json body from the request
		log.Print("Decoding JSON body")
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "parsing JSON body error", http.StatusBadRequest)
		}

		log.Print("Checking fields")
		if user.Username == "" || user.Password == "" {
			http.Error(w, "no username or password provided", http.StatusBadRequest)
		}

		_, err = db.AuthenticateCredentials(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		log.Print("Setting token expiration time")
		// Add expiration time for the claims
		expirationTime := time.Now().Add(time.Minute * 5)

		log.Print("Creating Claims")
		claims := &Claims{
			Username: user.Username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		log.Print("Getting token string")
		// Get new claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)

		if err != nil {
			http.Error(w, "unable to retrieve token", http.StatusInternalServerError)
			return
		}

		log.Print("Setting cookie")
		http.SetCookie(w,
			&http.Cookie{
				Name:    "token",
				Value:   tokenString,
				Expires: expirationTime,
			})
		log.Print("Successfully created cookie")
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (db *Database) SignUp(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		log.Print("Creating variables")
		var user User
		var creds map[string]string

		log.Print("Decoding JSON body")
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Print("Successfully decoded JSON body")

		log.Print("Checking fields")
		// check if fields are empty
		if creds["username"] == "" || creds["password"] == "" {
			http.Error(w, "missing fields", http.StatusBadRequest)
			return
		}
		log.Print("Fields not empty")

		log.Print("Authenticating credentials")
		// check if credentials exists
		user.Username = creds["username"]
		user.Password = creds["password"]
		_, err := db.AuthenticateCredentials(user)
		if err == nil {
			http.Error(w, "username already exists", http.StatusConflict)
			return
		}
		log.Print("Credentials are valid")

		// add new user credentials to the database
		db.Mu.Lock()
		defer db.Mu.Unlock()

		log.Print("Adding credentials to database")
		db.NextUserID++
		user.UserID = db.NextUserID
		db.Users = append(db.Users, user)
		log.Print("Successfully added credentials")

		w.Write([]byte("User created!"))
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (db *Database) LogOut(w http.ResponseWriter, r *http.Request) {

	// Destroy cookie
	c := http.Cookie{
		Name:   "token",
		MaxAge: -1}
	http.SetCookie(w, &c)

	w.Write([]byte("Logged out!\n"))
}

///////////////////////////////////////////////////////////////
// AUTHENTICATION
///////////////////////////////////////////////////////////////

func (db *Database) AuthenticateCredentials(user User) (int, error) {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	log.Print("Checking database for credentials")
	// Loop through database to find User
	for _, creds := range db.Users {
		if creds.Username == user.Username && creds.Password == user.Password {
			log.Print("Credentials found")
			return creds.UserID, nil
		}
	}
	log.Print("Credentials do not exist")
	return 0, errors.New("invalid username or password")
}

func (db *Database) AuthenticateRequest(w http.ResponseWriter, r *http.Request) error {
	// Get cookie from the request
	cookie, err := r.Cookie("token")
	if err != nil {
		if err != http.ErrNoCookie {
			http.Error(w, "no cookie found", http.StatusUnauthorized)
			return err
		}
		http.Error(w, "bad request", http.StatusBadRequest)
		return err
	}

	tokenStr := cookie.Value
	claims := &Claims{}

	// Get token
	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return err
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	// Check if token is valid
	if !tkn.Valid {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return errors.New("invalid token")
	}

	return nil
}

///////////////////////////////////////////////////////////////
// DATABASE OPERATIONS
///////////////////////////////////////////////////////////////
