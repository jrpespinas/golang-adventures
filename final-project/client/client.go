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
			Login(c)
			break choiceLoop
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

}

func SignUp(c http.Client) {
	// Print logo
	fmt.Printf("%sSign Up\n\n", logo)

	// Input credentials
	var username, password, confirmPassword string
	fmt.Print("Username: ")
	fmt.Scan(&username)
	fmt.Print("\nPassword: ")
	fmt.Scan(&password)
	fmt.Print("\nConfirm password: ")
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

	// POST request
	signupUrl := baseUrl + "/signup"
	request, err := http.NewRequest("POST", signupUrl, jsonBody)
	if err != nil {
		log.Printf("[SignUp] Error found: %s", err.Error())
		fmt.Printf("Error found: %s", err.Error())
		return
	}

	response, err := c.Do(request)
	if err != nil {
		fmt.Printf("Error found: %s", err.Error())
	}
	defer response.Body.Close()

	responseBody, _ := ioutil.ReadAll(response.Body)

	fmt.Printf("\nStatus: %s", response.Status)
	fmt.Printf("\nBody: %s", string(responseBody))
}

func Login(c http.Client) int {
	panic("Not implemented")
}
