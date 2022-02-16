package model

type Book struct {
	UserID int    `json:"userid"`
	BookID int    `json:"bookid"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Status string `json:"status"`
}
