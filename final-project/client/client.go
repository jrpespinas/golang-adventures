package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const baseUrl = "http://localhost:8080"

var logo = fmt.Sprintf("\n\nBook list\n=========\n\n")
var cookie []*http.Cookie

func main() {
	// initialize client
	c := http.Client{Timeout: time.Duration(1) * time.Second}

	// Display options
	welcome := fmt.Sprintf("%sChoose a command:\n(1) Login\n(2) Signup\n(3) Exit\n", logo)
	fmt.Println(welcome)

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
			}
			break
		case 2:
			SignUp(c)
		case 3:
			fmt.Printf("\n=========\nThank you for using Book list.")
			os.Exit(0)
		default:
			fmt.Printf("\nPlease choose a number.")
		}
		fmt.Println(welcome)
		fmt.Print("Option: ")
		fmt.Scan(&option)
	}

	//////////////////////////////////////////////
	fmt.Printf("%sChoose a command\n", logo)
	commands := []string{"View all books", "View finished books", "View unfinished books", "Find book", "Add book", "Edit book", "Delete book", "Logout"}
	for index, com := range commands {
		fmt.Printf("(%v) %v\n", index+1, com)
	}
	var bookOption int
	fmt.Print("Option: ")
	fmt.Scan(&bookOption)
	fmt.Println(bookOption)

secondLoop:
	for {
		switch bookOption {
		case 1:
			GetAllBooks(c)
		case 8:
			Logout(c)
			break secondLoop
		default:
			fmt.Printf("\nPlease choose a number.")
		}
	}

}

func SignUp(c http.Client) {
	// Print logo
	fmt.Printf("%sSign Up\n", logo)

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
	MakeRequest(c, "POST", "/signup", body, false)
	return
}

func Login(c http.Client) int {
	// Print logo
	fmt.Printf("%sLogin\n", logo)

	// Input credentials
	var username, password string
	fmt.Print("Username: ")
	fmt.Scan(&username)
	fmt.Print("Password: ")
	fmt.Scan(&password)

	// Encode json to byte
	body := fmt.Sprintf("{\"username\":\"%s\", \"password\":\"%s\"}", username, password)
	status, _ := MakeRequest(c, "POST", "/login", body, false)
	return status
}

func Logout(c http.Client) int {
	// GET request
	status, _ := MakeRequest(c, "GET", "/logout", "", true)
	return status
}

func GetAllBooks(c http.Client) {
	panic("Not implemented")
}

func MakeRequest(c http.Client, method string, endpoint string, body string, needCookie bool) (int, []byte) {
	// Set URL
	url := baseUrl + endpoint
	var request *http.Request

	// Initialize new request
	if method == "POST" || method == "PUT" {
		// Get JSON body
		jsonBody := bytes.NewBuffer([]byte(body))

		request, _ = http.NewRequest(method, url, jsonBody)
	} else if method == "GET" || method == "DELETE" {
		request, _ = http.NewRequest(method, url, nil)
	} else {
		return http.StatusMethodNotAllowed, []byte("null")
	}

	// Modify header
	request.Header.Set("Content-Type", "application/json")

	// Add cookie
	if needCookie {
		request.AddCookie(cookie[0])
	}

	response, err := c.Do(request)
	if err != nil {
		fmt.Printf("Error found: %s", err.Error())
	}
	defer response.Body.Close()

	responseBody, _ := ioutil.ReadAll(response.Body)

	cookie = response.Cookies()
	return response.StatusCode, responseBody
}
