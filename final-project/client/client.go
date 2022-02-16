package main

import (
	"bufio"
	"bytes"
	models "client/models"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const baseUrl = "http://localhost:8080"

var logo = fmt.Sprintf("\n\nBook list\n=========\n\n")
var cookie []*http.Cookie
var state = 1

func main() {
	// initialize client
	c := http.Client{Timeout: time.Duration(1) * time.Second}

	// Display options
	welcome := fmt.Sprintf("Choose a command:\n(1) Login\n(2) Signup\n(3) Exit\n")
	fmt.Print(logo)
	fmt.Print(welcome)

	// Choosing options
	var option int
	fmt.Print("Option: ")
	fmt.Scan(&option)

choiceLoop:
	for {
		switch option {
		case 1:
			status := Login(c)
			if status == 200 {
				break choiceLoop
			} else {
				fmt.Print("\nInvalid username and password")
			}
			break
		case 2:
			status := SignUp(c)
			if status == 200 {
				fmt.Println("\nUser created!")
			} else {
				fmt.Printf("\n%v\n", http.StatusText(status))
			}
		case 3:
			fmt.Printf("\n=========\nThank you for using Book list.")
			os.Exit(0)
		default:
			fmt.Printf("\nPlease choose a number.")
		}
		fmt.Print(logo)
		fmt.Println(welcome)
		fmt.Print("Option: ")
		fmt.Scan(&option)
	}

	//////////////////////////////////////////////
	if err := GetAllBooks(c); err != nil {
		fmt.Println("No books found")
	}
	fmt.Println("\nChoose a command:")
	commands := []string{"View finished books", "View unfinished books", "Find book", "Add book", "Edit book", "Delete book", "Logout"}
	for index, com := range commands {
		fmt.Printf("(%v) %v\n", index+1, com)
	}
	var bookOption int
	fmt.Print("Option: ")
	fmt.Scan(&bookOption)
secondLoop:
	for {
		switch bookOption {
		case 1:
			state = 2
		case 2:
			state = 3
		case 3:
			if err := GetOneBook(c); err != nil {
				fmt.Print("Book not found")
			}
		case 4:
			if err := AddBook(c); err != nil {
				fmt.Print(err.Error())
			}
		case 5:
			if err := EditBook(c); err != nil {
				fmt.Print(err.Error())
			}
		case 6:
			if err := DeleteBook(c); err != nil {
				fmt.Print(err.Error())
			}
		case 7:
			status := Logout(c)
			if status == 200 {
				fmt.Printf("%sThank you for using Book list", logo)
			}
			break secondLoop
		default:
			fmt.Printf("\nPlease choose a number.")
		}

		if state == 1 {
			if err := GetAllBooks(c); err != nil {
				fmt.Println("No books found")
			}
		} else if state == 2 {
			if err := GetFinishedBooks(c); err != nil {
				fmt.Println("No books found")
			}
		} else if state == 3 {
			if err := GetUnfinishedBooks(c); err != nil {
				fmt.Println("No books found")
			}
		}

		// Command prompt
		fmt.Println("\nChoose a command:")
		commands := []string{"View finished books", "View unfinished books", "Find book", "Add book", "Edit book", "Delete book", "Logout"}
		for index, com := range commands {
			fmt.Printf("(%v) %v\n", index+1, com)
		}
		fmt.Print("Option: ")
		fmt.Scan(&bookOption)
	}

}

func SignUp(c http.Client) int {
	// Print logo
	fmt.Print(logo)
	fmt.Println("Sign Up")

	// Input credentials
	var username, password, confirmPassword string
	fmt.Print("Username: ")
	fmt.Scan(&username)
	fmt.Print("Password: ")
	fmt.Scan(&password)
	fmt.Print("Confirm password: ")
	fmt.Scan(&confirmPassword)

	if password != confirmPassword {
		for confirmPassword != password {
			log.Print("[SignUp] Password did not match")
			log.Print("[SignUp] Re-entering passwords")
			fmt.Print("\nPassword: ")
			fmt.Scan(&password)
			fmt.Print("\nConfirm password: ")
			fmt.Scan(&confirmPassword)
		}
	}

	// Encode json to byte
	body := fmt.Sprintf("{\"username\":\"%s\", \"password\":\"%s\"}", username, password)
	jsonBody := bytes.NewBuffer([]byte(body))
	url := baseUrl + "/signup"

	// Prepare request
	req, err := http.NewRequest("POST", url, jsonBody)
	req.Header.Set("Content-Type", "application/json")

	// Make request
	res, err := c.Do(req)
	if err != nil {
		fmt.Printf("Error found: %s", err.Error())
	}
	defer res.Body.Close()

	return res.StatusCode
}

func Login(c http.Client) int {
	// Print logo
	fmt.Print(logo)
	fmt.Println("Login")

	// Input credentials
	var username, password string
	fmt.Print("Username: ")
	fmt.Scan(&username)
	fmt.Print("Password: ")
	fmt.Scan(&password)

	// Encode json to byte
	body := fmt.Sprintf("{\"username\":\"%s\", \"password\":\"%s\"}", username, password)
	jsonBody := bytes.NewBuffer([]byte(body))
	url := baseUrl + "/login"

	// Prepare request
	req, err := http.NewRequest("POST", url, jsonBody)
	req.Header.Set("Content-Type", "application/json")

	// Make request
	res, err := c.Do(req)
	if err != nil {
		fmt.Printf("Error found: %s", err.Error())
	}
	defer res.Body.Close()

	// Get cookie
	cookie = res.Cookies()

	// Return status
	return res.StatusCode
}

func Logout(c http.Client) int {
	// GET request
	// Encode json to byte
	url := baseUrl + "/logout"

	// Prepare request
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie[0])

	// Make request
	res, err := c.Do(req)
	if err != nil {
		fmt.Printf("Error found: %s", err.Error())
	}
	defer res.Body.Close()
	return res.StatusCode
}

func GetAllBooks(c http.Client) error {
	fmt.Print(logo)

	// Get URL
	url := baseUrl + "/books"

	// Prepare request
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie[0])

	// Make request
	res, err := c.Do(req)
	if err != nil {
		fmt.Printf("Error found: %s", err.Error())
	}
	defer res.Body.Close()

	// Decode JSON body
	var books []models.Book
	if err := json.NewDecoder(res.Body).Decode(&books); err != nil {
		return errors.New("No books found")
	}

	if len(books) != 0 {
		for index, book := range books {
			fmt.Printf("%v) Book ID: %v - %v by %v [%v]\n", index+1, book.BookID, book.Title, book.Author, book.Status)
		}
		return nil
	}
	return nil
}

func GetFinishedBooks(c http.Client) error {
	fmt.Print(logo)

	// Get URL
	url := baseUrl + "/books/finished"

	// Prepare request
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie[0])

	// Make request
	res, err := c.Do(req)
	if err != nil {
		fmt.Printf("Error found: %s", err.Error())
	}
	defer res.Body.Close()

	// Decode JSON body
	var books []models.Book
	if err := json.NewDecoder(res.Body).Decode(&books); err != nil {
		return errors.New("No books found")
	}

	if len(books) != 0 {
		for index, book := range books {
			fmt.Printf("%v) Book ID: %v - %v by %v [%v]\n", index+1, book.BookID, book.Title, book.Author, book.Status)
		}
		return nil
	}
	return nil
}

func GetUnfinishedBooks(c http.Client) error {
	fmt.Print(logo)

	// Get URL
	url := baseUrl + "/books/unfinished"

	// Prepare request
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie[0])

	// Make request
	res, err := c.Do(req)
	if err != nil {
		fmt.Printf("Error found: %s", err.Error())
	}
	defer res.Body.Close()

	// Decode JSON body
	var books []models.Book
	if err := json.NewDecoder(res.Body).Decode(&books); err != nil {
		return errors.New("No books found")
	}

	if len(books) != 0 {
		for index, book := range books {
			fmt.Printf("%v) Book ID: %v - %v by %v [%v]\n", index+1, book.BookID, book.Title, book.Author, book.Status)
		}
		return nil
	}
	return nil
}

func GetOneBook(c http.Client) error {
	// Print logo
	fmt.Print(logo)
	fmt.Println("Find Book")

	// Get Book ID
	var bookid int
	fmt.Print("Book ID: ")
	fmt.Scan(&bookid)

	// Get URL
	endpoint := "/books/" + strconv.Itoa(bookid)
	url := baseUrl + endpoint

	// Prepare request
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie[0])

	// Make request
	res, err := c.Do(req)
	if err != nil {
		fmt.Printf("Error found: %s", err.Error())
	}
	defer res.Body.Close()

	// Decode JSON body
	var book models.Book
	if err := json.NewDecoder(res.Body).Decode(&book); err != nil {
		return err
	}

	fmt.Printf("\nBook Found!\nBook ID: %v- %v by %v [%v]", book.BookID, book.Title, book.Author, book.Status)
	return nil
}

func AddBook(c http.Client) error {
	// Print logo
	fmt.Print(logo)
	fmt.Println("Add Book")

	// Input books
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Title: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSuffix(title, "\n")
	fmt.Print("Author: ")
	author, _ := reader.ReadString('\n')
	author = strings.TrimSuffix(author, "\n")
	fmt.Print("status: ")
	status, _ := reader.ReadString('\n')
	status = strings.TrimSuffix(status, "\n")

	// Encode json to byte
	body := fmt.Sprintf("{\"title\":\"%s\", \"author\":\"%s\", \"status\":\"%s\"}", title, author, status)
	jsonBody := bytes.NewBuffer([]byte(body))
	url := baseUrl + "/books"

	// Prepare request
	req, err := http.NewRequest("POST", url, jsonBody)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie[0])

	// Make request
	res, err := c.Do(req)
	if err != nil {
		fmt.Printf("Error found: %s", err.Error())
	}
	defer res.Body.Close()
	if res.StatusCode == 400 {
		return errors.New("Missing fields")
	}

	fmt.Printf("\n\nAdded %s to list!\n", title)
	return nil
}

func EditBook(c http.Client) error {
	// Print logo
	fmt.Print(logo)
	fmt.Println("Edit Book")

	// Get Book ID
	var bookid int
	fmt.Print("Book ID: ")
	fmt.Scan(&bookid)

	// Get URL
	endpoint := "/books/" + strconv.Itoa(bookid)
	url := baseUrl + endpoint

	// Edit book
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Title: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSuffix(title, "\n")
	fmt.Print("Author: ")
	author, _ := reader.ReadString('\n')
	author = strings.TrimSuffix(author, "\n")
	fmt.Print("status: ")
	status, _ := reader.ReadString('\n')
	status = strings.TrimSuffix(status, "\n")

	// Encode json to byte
	body := fmt.Sprintf("{\"title\":\"%s\", \"author\":\"%s\", \"status\":\"%s\"}", title, author, status)
	jsonBody := bytes.NewBuffer([]byte(body))

	// Prepare request
	req, err := http.NewRequest("PUT", url, jsonBody)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie[0])

	// Make request
	res, err := c.Do(req)
	if err != nil {
		fmt.Printf("Error found: %s", err.Error())
	}
	defer res.Body.Close()
	if res.StatusCode == 400 {
		return errors.New("Missing fields")
	}

	fmt.Printf("\n\nAdded %s to list!\n", title)
	return nil
}

func DeleteBook(c http.Client) error {
	// Print logo
	fmt.Print(logo)
	fmt.Println("Delete Book")

	// Get Book ID
	var bookid int
	fmt.Print("Book ID: ")
	fmt.Scan(&bookid)

	// Get URL
	endpoint := "/books/" + strconv.Itoa(bookid)
	url := baseUrl + endpoint

	// Prepare request
	req, err := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie[0])

	// Make request
	res, err := c.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode == 404 {
		return errors.New("Book not found")
	}
	defer res.Body.Close()
	fmt.Printf("\n\nDeleted Book %v\n", bookid)
	return nil
}
