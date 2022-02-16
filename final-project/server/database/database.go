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
		} else if r.URL.Path == "/books/finished" {
			db.GetAllRead(w, r)
		} else if r.URL.Path == "/books/unfinished" {
			db.GetAllUnread(w, r)
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
		log.Printf("[Login] %v Creating variables", r.Method)
		var user User

		// Get json body from the request
		log.Print("[Login] Decoding JSON body")
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "parsing JSON body error", http.StatusBadRequest)
		}

		log.Print("[Login] Checking fields")
		if user.Username == "" || user.Password == "" {
			http.Error(w, "no username or password provided", http.StatusBadRequest)
		}

		userID, err := db.AuthenticateCredentials(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		log.Print("[Login] Setting token expiration time")
		// Add expiration time for the claims
		expirationTime := time.Now().Add(time.Minute * 15)

		log.Print("[Login] Creating Claims")
		claims := &Claims{
			Username: user.Username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		log.Print("[Login] Getting token string")
		// Get new claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)

		if err != nil {
			http.Error(w, "unable to retrieve token", http.StatusInternalServerError)
			return
		}

		// Set cookie
		log.Print("[Login] Setting cookie")
		http.SetCookie(w,
			&http.Cookie{
				Name:    "token",
				Value:   tokenString,
				Expires: expirationTime,
			})

		// Track session
		log.Print("[Login] Creating session")
		var sess Session
		sess.Token = tokenString
		sess.UserID = userID
		db.Sessions = append(db.Sessions, sess)
		log.Print("[Login] Session created")

		w.Write([]byte("Login successful"))

	} else {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (db *Database) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		log.Print("[SignUp] Creating variables")
		var user User
		var creds map[string]string

		log.Print("[SignUp] Decoding JSON body")
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Print("[SignUp] Successfully decoded JSON body")

		log.Print("[SignUp] Checking fields")
		// check if fields are empty
		if creds["username"] == "" || creds["password"] == "" {
			http.Error(w, "missing fields", http.StatusBadRequest)
			return
		}
		log.Print("[SignUp] Fields not empty")

		log.Print("[SignUp] Authenticating credentials")
		// check if credentials exists
		user.Username = creds["username"]
		user.Password = creds["password"]
		_, err := db.AuthenticateCredentials(user)
		if err == nil {
			http.Error(w, "username already exists", http.StatusConflict)
			return
		}
		log.Print("[SignUp] Credentials are valid")

		// add new user credentials to the database
		db.Mu.Lock()

		log.Print("[SignUp] Adding credentials to database")
		db.NextUserID++
		user.UserID = db.NextUserID
		db.Users = append(db.Users, user)

		db.Mu.Unlock()

		log.Print("[SignUp] Successfully added credentials")

		w.Write([]byte("User created!"))
	} else {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (db *Database) LogOut(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Get token
		log.Print("[LogOut] Getting current session")
		_, token, err := db.GetSession(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Print("[LogOut] Deleting session")
		err = db.DeleteSession(token, w, r)
		if err != nil {
			log.Fatal("Session does not exist")
			http.Error(w, "Session does not exist", http.StatusBadRequest)
			return
		}

		// Destroy cookie
		log.Print("[LogOut] Destroying cookie")
		c := http.Cookie{
			Name:   "token",
			MaxAge: -1}
		http.SetCookie(w, &c)

		log.Print("[LogOut] Log out successful")
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

	log.Print("[AuthenticateCredentials] Checking database for credentials")
	// Loop through database to find User
	for _, creds := range db.Users {
		if creds.Username == user.Username && creds.Password == user.Password {
			log.Print("[AuthenticateCredentials] Credentials found")
			return creds.UserID, nil
		}
	}
	log.Print("[AuthenticateCredentials] Credentials do not exist")
	return 0, errors.New("user does not exist")
}

func (db *Database) AuthenticateRequest(w http.ResponseWriter, r *http.Request) (int, error) {
	// Get cookie from the request
	log.Print("[AuthenticateRequest] Getting cookie from headers")
	cookie, err := r.Cookie("token")
	if err != nil {
		if err != http.ErrNoCookie {
			log.Printf("[AuthenticateRequest] Error found: %s", err.Error())
			return http.StatusBadRequest, err
		}
		log.Printf("[AuthenticateRequest] Error found: %s", err.Error())
		return http.StatusBadRequest, err
	}

	tokenStr := cookie.Value
	claims := &Claims{}

	// Get token
	log.Print("[AuthenticateRequest] Parsing claims")
	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			log.Printf("[AuthenticateRequest] Error found: %s", err.Error())
			return http.StatusUnauthorized, err
		}
		log.Printf("[AuthenticateRequest] Error found: %s", err.Error())
		return http.StatusBadRequest, err
	}

	// Check if token is valid
	log.Print("[AuthenticateRequest] Validating token")
	if !tkn.Valid {
		log.Print("[AuthenticateRequest] Invalid token")
		return http.StatusUnauthorized, errors.New("invalid token")
	}

	log.Print("[AuthenticateRequest] Request authenticated")
	return http.StatusOK, nil
}

func (db *Database) GetSession(w http.ResponseWriter, r *http.Request) (int, string, error) {
	log.Print("[GetSession] Getting session token")
	cookie, err := r.Cookie("token")
	if err != nil {
		if err != http.ErrNoCookie {
			log.Printf("[GetSession] Error found: %s", err.Error())
			return 0, "", err
		}
		log.Print(err.Error())
		return 0, "", err
	}

	log.Print("[GetSession] Checking database for credentials")
	// Loop through database to find User
	for _, creds := range db.Sessions {
		if creds.Token == cookie.Value {
			log.Print("[GetSession] Session found")
			return creds.UserID, cookie.Value, nil
		}
	}
	msg := "session does not exist"
	log.Printf("[GetSession] Error found: %s", msg)
	return 0, "", errors.New(msg)
}

func (db *Database) DeleteSession(token string, w http.ResponseWriter, r *http.Request) error {
	db.Mu.Lock()
	defer db.Mu.Unlock()
	for index, sess := range db.Sessions {
		if sess.Token == token {
			db.Sessions = append(db.Sessions[:index], db.Sessions[index+1:]...)
			log.Print("[DeleteSession] Successfully delete session")
			return nil
		}
	}

	return errors.New("Session not found")
}

///////////////////////////////////////////////////////////////
// DATABASE OPERATIONS
///////////////////////////////////////////////////////////////

func (db *Database) Update(w http.ResponseWriter, r *http.Request) {
	log.Print("[Update] Updating database")

	db.Mu.Lock()

	byteData, err := json.MarshalIndent(db, "", "	")
	if err != nil {
		log.Fatalf("[Update] Marshaling of data failed: %s\n", err.Error())
	}

	log.Print("[Update] Writing to JSON file")
	if err := ioutil.WriteFile("data/data.json", byteData, 0644); err != nil {
		log.Fatalf("[Update] Failed to update database: %s\n", err.Error())
	}
	db.Mu.Unlock()
}

///////////////////////////////////////////////////////////////
// BOOK OPERATIONS
///////////////////////////////////////////////////////////////

func (db *Database) GetAllRead(w http.ResponseWriter, r *http.Request) {
	// Authenticate request
	log.Print("[GetAllRead] Authenticate request")
	status, err := db.AuthenticateRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	// Get token and userID
	log.Print("[GetAllRead] Get session")
	userID, _, err := db.GetSession(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	switch r.Method {
	case "POST":
		// Create placeholder for book
		var book map[string]string

		// Decode json body
		log.Print("[GetAllRead] Decoding JSON body")
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			log.Print("[GetAllRead] Decoding JSON body failed")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Check for missing fields
		log.Print("[GetAllRead] Check fields")
		if book["title"] == "" || book["author"] == "" || book["status"] == "" {
			log.Print("[GetAllRead] Error found: missing fields")
			http.Error(w, "missing fields", http.StatusBadRequest)
			return
		}

		// Check database for similar entries
		log.Print("[GetAllRead] Checking database for similar entry")
		var newbook Book
		newbook.Title = book["title"]
		newbook.Author = book["author"]
		newbook.Status = book["status"]
		newbook.UserID = userID

		err = db.CheckBook(newbook)
		if err == nil {
			log.Print("[GetAllRead] Error found: book already exists")
			http.Error(w, "book already exists", http.StatusConflict)
			return
		}

		// Add book to database
		log.Print("[GetAllRead] Adding Book to database")
		db.Mu.Lock()

		db.NextBookID++
		newbook.BookID = db.NextBookID
		db.Books = append(db.Books, newbook)
		db.Mu.Unlock()

		log.Print("[GetAllRead] Successfully added book to database")
		w.Write([]byte("Book added"))

	case "GET":
		books, status, err := db.GetBookByUser(userID, "finished", w, r)
		if err != nil {
			w.Write([]byte(err.Error()))
			http.Error(w, err.Error(), status)
			return
		}

		// Return all books
		log.Print("[GetAllRead] Encoding")
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(books); err != nil {
			log.Printf("[GetAllRead] Error found: %s", err.Error())
			w.Write([]byte(err.Error()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Print("[GetAllRead] JSON encoded")

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func (db *Database) GetAllUnread(w http.ResponseWriter, r *http.Request) {
	// Authenticate request
	log.Print("[GetAllUnead] Authenticate request")
	status, err := db.AuthenticateRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	// Get token and userID
	log.Print("[GetAllUnead] Get session")
	userID, _, err := db.GetSession(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	switch r.Method {
	case "POST":
		// Create placeholder for book
		var book map[string]string

		// Decode json body
		log.Print("[GetAllUnead] Decoding JSON body")
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			log.Printf("[GetAllUnead] Error found: %s", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Check for missing fields
		log.Print("[GetAllUnead] Check fields")
		if book["title"] == "" || book["author"] == "" || book["status"] == "" {
			log.Print("[GetAllUnead] Error found: missing fields")
			http.Error(w, "missing fields", http.StatusBadRequest)
			return
		}

		// Check database for similar entries
		log.Print("[GetAllUnead] Checking database for similar entry")
		var newbook Book
		newbook.Title = book["title"]
		newbook.Author = book["author"]
		newbook.Status = book["status"]
		newbook.UserID = userID

		err = db.CheckBook(newbook)
		if err == nil {
			log.Printf("[GetAllUnead] Error found: %s", err.Error())
			http.Error(w, "book already exists", http.StatusConflict)
			return
		}

		// Add book to database
		log.Print("[GetAllUnead] Adding Book to database")
		db.Mu.Lock()

		db.NextBookID++
		newbook.BookID = db.NextBookID
		db.Books = append(db.Books, newbook)
		db.Mu.Unlock()

		log.Print("[GetAllUnead] Successfully added book to database")
		w.Write([]byte("Book added"))

	case "GET":
		books, status, err := db.GetBookByUser(userID, "not", w, r)
		if err != nil {
			http.Error(w, err.Error(), status)
			return
		}

		// Return all books
		log.Print("[GetAllUnead] Encoding")
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(books); err != nil {
			log.Printf("[GetAllUnead] Error found: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Print("[GetAllUnead] JSON encoded")

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func (db *Database) ProcessBooksID(bookID int, w http.ResponseWriter, r *http.Request) {
	// Authenticate request
	log.Print("[ProcessBooksID] Authenticate request")
	status, err := db.AuthenticateRequest(w, r)
	if err != nil {
		log.Printf("[ProcessBooksID] Error found: %s", err.Error())
		http.Error(w, err.Error(), status)
		return
	}

	// Get token and userID
	log.Print("[ProcessBooksID] Get session")
	userID, _, err := db.GetSession(w, r)
	if err != nil {
		log.Printf("[ProcessBooksID] Error found: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// return one book
	book, err := db.GetBookByID(userID, bookID)
	if err != nil {
		log.Printf("[ProcessBooksID] Error found: %s", err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	log.Printf("[ProcessBooksID] Book %d found", book.BookID)

	switch r.Method {
	case "POST":
		w.WriteHeader(http.StatusMethodNotAllowed)
	case "GET":
		log.Print("[ProcessBooksID] Encoding book")
		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(book); err != nil {
			log.Printf("[ProcessBooksID] Error found: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Print("[ProcessBooksID] JSON encoded")
	case "PUT":
		var bookPlaceholder map[string]string
		// Decode JSON body
		log.Print("[ProcessBooksID] Decoding JSON body")
		if err := json.NewDecoder(r.Body).Decode(&bookPlaceholder); err != nil {
			log.Printf("[ProcessBooksID] Error found: %s", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Check for missing fields
		log.Print("[ProcessBooksID] Check fields")
		if bookPlaceholder["title"] == "" || bookPlaceholder["author"] == "" || bookPlaceholder["status"] == "" {
			log.Print("[ProcessBooksID] Error found: missing fields")
			http.Error(w, "missing fields", http.StatusBadRequest)
			return
		}

		// Delete book to be replaced
		err = db.DeleteBookByID(userID, bookID)
		// db.Update(w, r)
		if err != nil {
			log.Printf("[ProcessBooksID] Error found: %s", err.Error())
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		log.Print("[ProcessBooksID] Creating book placeholder")
		var newbook Book
		newbook.Title = bookPlaceholder["title"]
		newbook.Author = bookPlaceholder["author"]
		newbook.Status = bookPlaceholder["status"]
		newbook.BookID = book.BookID
		newbook.UserID = book.UserID

		log.Print("[ProcessBooksID] Updating book information")
		db.Mu.Lock()
		db.Books = append(db.Books, newbook)
		db.Mu.Unlock()

		log.Print("[ProcessBooksID] Successfully updated book information")
		w.Write([]byte("[ProcessBooksID] Book updated"))
	case "DELETE":
		// Delete book to be replaced
		if err = db.DeleteBookByID(userID, bookID); err != nil {
			log.Printf("[ProcessBooksID] Error found: %s", err.Error())
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Write([]byte("Book deleted"))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (db *Database) ProcessBooks(w http.ResponseWriter, r *http.Request) {
	// Authenticate GET request
	log.Print("[ProcessBooks] Authenticate request")
	status, err := db.AuthenticateRequest(w, r)
	if err != nil {
		log.Printf("[ProcessBooks] Error found: %s", err.Error())
		http.Error(w, err.Error(), status)
		return
	}

	// Get token and userID
	log.Print("[ProcessBooks] Get session")
	userID, _, err := db.GetSession(w, r)
	if err != nil {
		log.Printf("[ProcessBooks] Error found: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	switch r.Method {
	case "POST":
		// Create placeholder for book
		var book map[string]string

		// Decode json body
		log.Print("[ProcessBooks] Decoding JSON body")
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			log.Printf("[ProcessBooks] Error found: %s", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Check for missing fields
		log.Print("[ProcessBooks] Check fields")
		if book["title"] == "" || book["author"] == "" || book["status"] == "" {
			log.Print("[ProcessBooks] Error found: missing fields")
			http.Error(w, "missing fields", http.StatusBadRequest)
			return
		}

		// Check database for similar entries
		log.Print("[ProcessBooks] Checking database for similar entry")
		var newbook Book
		newbook.Title = book["title"]
		newbook.Author = book["author"]
		newbook.Status = book["status"]
		newbook.UserID = userID

		if err = db.CheckBook(newbook); err == nil {
			log.Printf("[ProcessBooks] Error found: %s", err.Error())
			http.Error(w, "book already exists", http.StatusConflict)
			return
		}

		// Add book to database
		log.Print("[ProcessBooks] Adding Book to database")
		db.Mu.Lock()

		db.NextBookID++
		newbook.BookID = db.NextBookID
		db.Books = append(db.Books, newbook)
		db.Mu.Unlock()

		log.Print("[ProcessBooks] Successfully added book to database")
		w.Write([]byte("Book added"))
	case "GET":
		// Get books by userID
		books, status, err := db.GetBookByUser(userID, "all", w, r)
		if err != nil {
			http.Error(w, err.Error(), status)
			return
		}

		// Return all books
		log.Print("[ProcessBooks] Encoding")
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(books); err != nil {
			log.Printf("[ProcessBooks] Error found: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Print("[ProcessBooks] JSON encoded")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (db *Database) CheckBook(book Book) error {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	log.Printf("[CheckBook] Checking database for %s by %s", book.Title, book.Author)
	for _, bookEntry := range db.Books {
		if bookEntry.UserID == book.UserID && bookEntry.Author == book.Author && bookEntry.Title == book.Title {
			log.Print("[CheckBook] Book found")
			return nil
		}
	}
	return errors.New("Book already exists")
}

func (db *Database) GetBookByUser(userID int, category string, w http.ResponseWriter, r *http.Request) ([]Book, int, error) {
	booksByUser := []Book{}

	db.Mu.Lock()
	defer db.Mu.Unlock()

	if category == "finished" {
		log.Print("[GetBookByUser] Getting all finished books by user")
		for _, book := range db.Books {
			if book.UserID == userID && book.Status == "finished" {
				booksByUser = append(booksByUser, book)
			}
		}
	} else if category == "not" {
		log.Print("[GetBookByUser] Getting all unfinished books by user")
		for _, book := range db.Books {
			if book.UserID == userID && book.Status != "finished" {
				booksByUser = append(booksByUser, book)
			}
		}
	} else if category == "all" {
		log.Print("[GetBookByUser] Getting all books by user")
		for _, book := range db.Books {
			if book.UserID == userID {
				booksByUser = append(booksByUser, book)
			}
		}
	} else {
		log.Print("[GetBookByUser] Category not recognized")
		return booksByUser, http.StatusBadRequest, errors.New("category not recognized")
	}

	if len(booksByUser) == 0 {
		log.Print("[GetBookByUser] No books found")
		return booksByUser, http.StatusNotFound, errors.New("no books found")
	}
	return booksByUser, http.StatusOK, nil
}

func (db *Database) GetBookByID(userID int, bookID int) (Book, error) {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	log.Printf("[GetBookByID] Getting Book ID: %d", bookID)
	for _, book := range db.Books {
		if book.UserID == userID && book.BookID == bookID {
			return book, nil
		}
	}
	return Book{}, errors.New("Book does not exist")
}

func (db *Database) DeleteBookByID(userID int, bookID int) error {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	log.Printf("[DeleteBookByID] Deleting Book ID: %d", bookID)
	for index, book := range db.Books {
		if book.UserID == userID && book.BookID == bookID {
			db.Books = append(db.Books[:index], db.Books[index+1:]...)
			log.Print("[DeleteBookByID] Successfully deleted book")
			return nil
		}
	}
	return errors.New("Book not found")
}
