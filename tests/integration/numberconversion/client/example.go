package numberconversion

// This file contains usage examples for the generated SOAP client
// To use this client in your code:
//
// import "your-module/tests/integration/numberconversion/client"
//
// Example usage:

/*
package main

import (
	"fmt"
	"log"

	"tests/integration/numberconversion/client"
)

func main() {
	// Create a new client
	client := numberconversion.NewClient("")

	// You can also specify a custom URL:
	// client := numberconversion.NewClient("http://your-service-url")

	// Example: Call NumberToWords operation
	result, err := client.NumberToWords(nil)
	if err != nil {
		log.Fatalf("Failed to call NumberToWords: %v", err)
	}

	fmt.Printf("Result: %+v\n", result)
}
*/

// Available Operations:
//
// client.NumberToWords(parameters ) (, error)
//   Returns the word corresponding to the positive number passed as parameter. Limited to quadrillions.
//
// client.NumberToDollars(parameters ) (, error)
//   Returns the non-zero dollar amount of the passed number.
//
