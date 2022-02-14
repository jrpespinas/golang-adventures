package models

import (
	"sync"

	"github.com/dgrijalva/jwt-go"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	UserID   int    `json:"userid"`
}

type Book struct {
	UserID int    `json:"userid"`
	BookID int    `json:"bookid"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Isbn   string `json:"isbn"`
	Status string `json:"status"`
}

type Database struct {
	Users      []User `json:"users"`
	Books      []Book `json:"books"`
	NextUserID int    `json:"nextUserID"`
	NextBookID int    `json:"nextBookID"`
	Mu         sync.Mutex
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
