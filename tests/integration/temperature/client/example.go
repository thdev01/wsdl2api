package temperature

// This file contains usage examples for the generated SOAP client
// To use this client in your code:
//
// import "your-module/tests/integration/temperature/client"
//
// Example usage:

/*
package main

import (
	"fmt"
	"log"

	"tests/integration/temperature/client"
)

func main() {
	// Create a new client
	client := temperature.NewClient("")

	// You can also specify a custom URL:
	// client := temperature.NewClient("http://your-service-url")

	// Example: Call CelsiusToFahrenheit operation
	result, err := client.CelsiusToFahrenheit(nil)
	if err != nil {
		log.Fatalf("Failed to call CelsiusToFahrenheit: %v", err)
	}

	fmt.Printf("Result: %+v\n", result)
}
*/

// Available Operations:
//
// client.CelsiusToFahrenheit(parameters ) (, error)
//   Converts a Celsius Temperature to a Fahrenheit value
//
// client.FahrenheitToCelsius(parameters ) (, error)
//   Converts a Fahrenheit Temperature to a Celsius value
//
// client.WindChillInCelsius(parameters ) (, error)
//   Windchill temperature calculated with the formula of Steadman
//
// client.WindChillInFahrenheit(parameters ) (, error)
//   Windchill temperature calculated with the formula of Steadman
//
