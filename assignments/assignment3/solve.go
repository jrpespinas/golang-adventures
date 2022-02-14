package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Log unexisting URL path
	if r.URL.Path != "/solve" {
		log.Fatalf("url does not exist: %v", r.URL.Path)
	}

	// Parse coefficients
	coefficients := getCoefficientsMatrix(w, r)

	// Print systems of equation
	fmt.Fprintln(w, "System:")
	for i, j := 0, 0; i < 3; i, j = i+1, j+4 {
		fmt.Fprintf(w, "%dx + %dy + %dz = %d\n", coefficients[j], coefficients[j+1], coefficients[j+2], coefficients[j+3])
	}

	// return solution
	d, dx, dy, dz := computeDeterminants(coefficients)

	if d != 0 {
		fmt.Fprintf(w, "\nSolution:\nx = %.2f, y = %.2f, z = %.2f\n", dx/d, dy/d, dz/d)
	} else if d == 0 && dx == 0 && dy == 0 && dz == 0 {
		fmt.Fprintf(w, "\nDependent - With Multiple Solutions")
	} else {
		fmt.Fprintf(w, "\nInconsistent - No Solution")
	}
}

func getCoefficientsMatrix(w http.ResponseWriter, r *http.Request) []int {
	if err := r.ParseForm(); err != nil {
		log.Fatalf("Can not parse form: %v", err)
	}

	// initialize coefficients array
	var coefficients []int

	if nums, ok := r.Form["coef"]; ok {
		coefStrigns := strings.Split(nums[0], ",")

		// check number of coefficients
		if len(coefStrigns) != 12 {
			log.Fatalln("Insufficient number of `coef`: provide 12 integers")
		}

		// convert coefficients to integers
		for _, coef := range coefStrigns {
			if val, err := strconv.Atoi(coef); err == nil {
				coefficients = append(coefficients, val)
			} else {
				log.Fatalf("parsed value not an integer: %v", err)
			}
		}
	}
	return coefficients
}

func computeDeterminants(coefficients []int) (float64, float64, float64, float64) {
	d := float64((coefficients[0] * (coefficients[5]*coefficients[10] - coefficients[6]*coefficients[9])) - (coefficients[4] * (coefficients[1]*coefficients[10] - coefficients[2]*coefficients[9])) + (coefficients[8] * (coefficients[1]*coefficients[6] - coefficients[2]*coefficients[5])))
	dx := float64((coefficients[3] * (coefficients[5]*coefficients[10] - coefficients[6]*coefficients[9])) - (coefficients[7] * (coefficients[1]*coefficients[10] - coefficients[2]*coefficients[9])) + (coefficients[11] * (coefficients[1]*coefficients[6] - coefficients[2]*coefficients[5])))
	dy := float64((coefficients[0] * (coefficients[7]*coefficients[10] - coefficients[6]*coefficients[11])) - (coefficients[4] * (coefficients[3]*coefficients[10] - coefficients[2]*coefficients[11])) + (coefficients[8] * (coefficients[3]*coefficients[6] - coefficients[2]*coefficients[7])))
	dz := float64((coefficients[0] * (coefficients[5]*coefficients[11] - coefficients[7]*coefficients[9])) - (coefficients[4] * (coefficients[1]*coefficients[11] - coefficients[3]*coefficients[9])) + (coefficients[8] * (coefficients[1]*coefficients[7] - coefficients[3]*coefficients[5])))

	return d, dx, dy, dz
}
