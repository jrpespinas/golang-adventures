package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("secret_key")

func (db *Database) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var bookID int
		if r.URL.Path == "/login" {
			db.Login(w, r)
		} else if r.URL.Path == "/signup" {
			db.SignUp(w, r)
		} else if r.URL.Path == "/logout" {
			db.LogOut(w, r)
		} else if r.URL.Path == "/books" {
			db.ProcessBooks(w, r)
		} else if n, _ := fmt.Sscanf(r.URL.Path, "/books/%d", &bookID); n == 1 {
			db.ProcessBooksID(bookID, w, r)
		} else {
			http.Error(w, "url does not exist", http.StatusNotImplemented)
		}
		db.Update(w, r)
	}
}

///////////////////////////////////////////////////////////////
// LOGGING FEATURES
///////////////////////////////////////////////////////////////

func (db *Database) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		log.Printf("%v Creating variables", r.Method)
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

		userID, err := db.AuthenticateCredentials(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		log.Print("Setting token expiration time")
		// Add expiration time for the claims
		expirationTime := time.Now().Add(time.Minute * 15)

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

		// Set cookie
		log.Print("Setting cookie")
		http.SetCookie(w,
			&http.Cookie{
				Name:    "token",
				Value:   tokenString,
				Expires: expirationTime,
			})
		log.Print("Successfully created cookie")

		// Track session
		log.Print("Creating session")
		var sess Session
		sess.Token = tokenString
		sess.UserID = userID
		db.Sessions = append(db.Sessions, sess)
		log.Print("Session created")

	} else {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (db *Database) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
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

		log.Print("Adding credentials to database")
		db.NextUserID++
		user.UserID = db.NextUserID
		db.Users = append(db.Users, user)

		db.Mu.Unlock()

		log.Print("Successfully added credentials")

		w.Write([]byte("User created!"))
	} else {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (db *Database) LogOut(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Get token
		log.Print("Getting current session")
		_, token, err := db.GetSession(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Print("Deleting session")
		err = db.DeleteSession(token, w, r)
		if err != nil {
			log.Fatal("Session does not exist")
			http.Error(w, "Session does not exist", http.StatusBadRequest)
			return
		}

		// Destroy cookie
		log.Print("Destroying cookie")
		c := http.Cookie{
			Name:   "token",
			MaxAge: -1}
		http.SetCookie(w, &c)

		log.Print("Log out successful")
		w.Write([]byte("Logged out!\n"))
	} else {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
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
	return 0, errors.New("user does not exist")
}

func (db *Database) AuthenticateRequest(w http.ResponseWriter, r *http.Request) (int, error) {
	// Get cookie from the request
	log.Print("Getting cookie from headers")
	cookie, err := r.Cookie("token")
	if err != nil {
		if err != http.ErrNoCookie {
			log.Print(err.Error())
			return http.StatusBadRequest, err
		}
		log.Print(err.Error())
		return http.StatusBadRequest, err
	}

	tokenStr := cookie.Value
	claims := &Claims{}

	// Get token
	log.Print("Parsing claims")
	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			log.Print(err.Error())
			return http.StatusUnauthorized, err
		}
		log.Print(err.Error())
		return http.StatusBadRequest, err
	}

	// Check if token is valid
	log.Print("Validating token")
	if !tkn.Valid {
		log.Print("Invalid token")
		return http.StatusUnauthorized, errors.New("invalid token")
	}

	log.Print("Request authenticated")
	return http.StatusOK, nil
}

func (db *Database) GetSession(w http.ResponseWriter, r *http.Request) (int, string, error) {
	log.Print("Getting session token")
	cookie, err := r.Cookie("token")
	if err != nil {
		if err != http.ErrNoCookie {
			log.Print(err.Error())
			return 0, "", err
		}
		log.Print(err.Error())
		return 0, "", err
	}

	log.Print("Checking database for credentials")
	// Loop through database to find User
	for _, creds := range db.Sessions {
		if creds.Token == cookie.Value {
			log.Print("Session found")
			return creds.UserID, cookie.Value, nil
		}
	}
	msg := "session does not exist"
	log.Print(msg)
	return 0, "", errors.New(msg)
}

func (db *Database) DeleteSession(token string, w http.ResponseWriter, r *http.Request) error {
	db.Mu.Lock()
	defer db.Mu.Unlock()
	for index, sess := range db.Sessions {
		if sess.Token == token {
			db.Sessions = append(db.Sessions[:index], db.Sessions[index+1:]...)
			log.Print("Successfully delete session")
			return nil
		}
	}

	return errors.New("Session not found")
}

///////////////////////////////////////////////////////////////
// DATABASE OPERATIONS
///////////////////////////////////////////////////////////////

func (db *Database) Update(w http.ResponseWriter, r *http.Request) {
	log.Print("Updating database")

	db.Mu.Lock()

	byteData, err := json.MarshalIndent(db, "", "	")
	if err != nil {
		log.Fatalf("Marshaling of data failed: %s\n", err.Error())
	}

	log.Print("Writing to JSON file")
	if err := ioutil.WriteFile("data/data.json", byteData, 0644); err != nil {
		log.Fatalf("Failed to update database: %s\n", err.Error())
	}
	db.Mu.Unlock()
}

///////////////////////////////////////////////////////////////
// BOOK OPERATIONS
///////////////////////////////////////////////////////////////

func (db *Database) ProcessBooksID(bookID int, w http.ResponseWriter, r *http.Request) {
	// Authenticate GET request
	log.Print("Authenticate request")
	status, err := db.AuthenticateRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	// Get token and userID
	log.Print("Get session")
	_, _, err = db.GetSession(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	switch r.Method {
	case "POST":
		w.WriteHeader(http.StatusMethodNotAllowed)
	case "GET":
		panic("Not implemented")
	case "PUT":
		panic("Not implemented")
	case "DELETE":
		panic("Not implemented")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (db *Database) ProcessBooks(w http.ResponseWriter, r *http.Request) {
	// Authenticate GET request
	log.Print("Authenticate request")
	status, err := db.AuthenticateRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	// Get token and userID
	log.Print("Get session")
	userID, _, err := db.GetSession(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	switch r.Method {
	case "POST":
		// Create placeholder for book
		var book map[string]string

		log.Print("Decoding JSON body")
		// Decode json body
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			log.Fatalf("Decoding JSON body failed")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Check for missing fields
		log.Print("Check fields")
		if book["title"] == "" || book["author"] == "" || book["status"] == "" {
			http.Error(w, "missing fields", http.StatusBadRequest)
			return
		}

		// Check database for similar entries
		log.Print("Checking database for similar entry")
		var newbook Book
		newbook.Title = book["title"]
		newbook.Author = book["author"]
		newbook.Status = book["status"]
		newbook.UserID = userID

		err = db.CheckBook(newbook)
		if err == nil {
			http.Error(w, "book already exists", http.StatusConflict)
			return
		}

		// Add book to database
		log.Print("Adding Book to database")
		db.Mu.Lock()

		db.NextBookID++
		newbook.BookID = db.NextBookID
		db.Books = append(db.Books, newbook)
		db.Mu.Unlock()

		log.Print("Successfully added book to database")
		w.Write([]byte("Book added"))
	case "GET":
		// Get books by userID
		books := db.GetBookByUser(userID)

		// Return all books
		log.Print("Encoding entries")
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(books); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (db *Database) CheckBook(book Book) error {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	log.Printf("Checking database for %s by %s", book.Title, book.Author)
	for _, bookEntry := range db.Books {
		if bookEntry.UserID == book.UserID && bookEntry.Author == book.Author && bookEntry.Title == book.Title && bookEntry.Status == book.Status {
			log.Print("Book found")
			return nil
		}
	}
	return errors.New("Book already exists")
}

func (db *Database) GetBookByUser(userID int) []Book {
	booksByUser := []Book{}

	db.Mu.Lock()
	defer db.Mu.Unlock()

	log.Print("Getting books by user")
	for _, book := range db.Books {
		if book.UserID == userID {
			booksByUser = append(booksByUser, book)
		}
	}
	return booksByUser
}
