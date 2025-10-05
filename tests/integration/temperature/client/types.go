package temperature

import "encoding/xml"

// Auto-generated types from WSDL

// CelsiusToFahrenheitRequest represents CelsiusToFahrenheit request
type CelsiusToFahrenheitRequest struct {
	XMLName  xml.Name `xml:"CelsiusToFahrenheit"`
	NCelsius string   `xml:"nCelsius"`
}

// CelsiusToFahrenheitResponse represents CelsiusToFahrenheit response
type CelsiusToFahrenheitResponse struct {
	XMLName                        xml.Name `xml:"CelsiusToFahrenheitResponse"`
	NCelsiusToFahrenheitResult string   `xml:"CelsiusToFahrenheitResult"`
}

// FahrenheitToCelsiusRequest represents FahrenheitToCelsius request
type FahrenheitToCelsiusRequest struct {
	XMLName     xml.Name `xml:"FahrenheitToCelsius"`
	NFahrenheit string   `xml:"nFahrenheit"`
}

// FahrenheitToCelsiusResponse represents FahrenheitToCelsius response
type FahrenheitToCelsiusResponse struct {
	XMLName                        xml.Name `xml:"FahrenheitToCelsiusResponse"`
	NFahrenheitToCelsiusResult string   `xml:"FahrenheitToCelsiusResult"`
}
