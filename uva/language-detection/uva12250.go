package main

import "fmt"

func main() {
	i := 1
	for true {
		var lang string
		fmt.Scan(&lang)
		if lang == "#" {
			break
		}

		switch lang {
		case "HELLO":
			lang = "ENGLISH"
		case "HOLA":
			lang = "SPANISH"
		case "HALLO":
			lang = "GERMAN"
		case "BONJOUR":
			lang = "FRENCH"
		case "CIAO":
			lang = "ITALIAN"
		case "ZDRAVSTVUJTE":
			lang = "RUSSIAN"
		default:
			lang = "UNKNOWN"
		}
		fmt.Printf("Case %v: %v\n", i, lang)
		i++
	}
}
