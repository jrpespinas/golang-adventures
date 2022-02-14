package models

import (
	"encoding/json"
	"errors"
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
		} else if r.URL.Path == "/Logout" {
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
	var user User

	// Get json body from the request
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "parsing JSON body error", http.StatusBadRequest)
	}

	if user.Username == "" || user.Password == "" {
		http.Error(w, "no username or password provided", http.StatusBadRequest)
	}

	_, err = db.AuthenticateCredentials(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Add expiration time for the claims
	expirationTime := time.Now().Add(time.Minute * 5)

	claims := &Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Get new claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		http.Error(w, "unable to retrieve token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w,
		&http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})
}

func (db *Database) SignUp(w http.ResponseWriter, r *http.Request) {
	var user User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// check if fields are empty
	if user.Username == "" || user.Password == "" {
		http.Error(w, "missing fields", http.StatusBadRequest)
	}

	db.Mu.Lock()
	defer db.Mu.Unlock()

	// check if credentials exists
	_, err := db.AuthenticateCredentials(user)
	if err == nil {
		http.Error(w, "username already exists", http.StatusConflict)
	}

	// add new user credentials to the database
	db.NextUserID++
	user.UserID = db.NextUserID
	db.Users = append(db.Users, user)

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

	// Loop through database to find User
	for _, creds := range db.Users {
		if creds.Username == user.Username && creds.Password == user.Password {
			return creds.UserID, nil
		}
	}
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
