package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
)

// CalculatorClient is a simple SOAP client for calculator service
type CalculatorClient struct {
	URL string
}

// NewCalculatorClient creates a new calculator client
func NewCalculatorClient() *CalculatorClient {
	return &CalculatorClient{
		URL: "http://www.dneonline.com/calculator.asmx",
	}
}

// AddRequest represents the Add operation request
type AddRequest struct {
	XMLName xml.Name `xml:"http://tempuri.org/ Add"`
	IntA    int      `xml:"intA"`
	IntB    int      `xml:"intB"`
}

// AddResponse represents the Add operation response
type AddResponse struct {
	XMLName   xml.Name `xml:"AddResponse"`
	AddResult int      `xml:"AddResult"`
}

// SOAPEnvelope represents a SOAP envelope
type SOAPEnvelope struct {
	XMLName xml.Name `xml:"soap:Envelope"`
	Soap    string   `xml:"xmlns:soap,attr"`
	Body    SOAPBody
}

// SOAPBody represents the SOAP body
type SOAPBody struct {
	XMLName xml.Name `xml:"soap:Body"`
	Content interface{}
}

// Add performs addition of two integers
func (c *CalculatorClient) Add(a, b int) (int, error) {
	// Create request
	request := AddRequest{
		IntA: a,
		IntB: b,
	}

	// Build SOAP envelope
	envelope := SOAPEnvelope{
		Soap: "http://schemas.xmlsoap.org/soap/envelope/",
		Body: SOAPBody{
			Content: request,
		},
	}

	// Marshal to XML
	xmlData, err := xml.MarshalIndent(envelope, "", "  ")
	if err != nil {
		return 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Add XML header
	requestBody := []byte(xml.Header + string(xmlData))

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", c.URL, bytes.NewReader(requestBody))
	if err != nil {
		return 0, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "text/xml; charset=utf-8")
	httpReq.Header.Set("SOAPAction", "\"http://tempuri.org/Add\"")

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return 0, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var responseEnvelope struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    struct {
			XMLName  xml.Name    `xml:"Body"`
			Response AddResponse `xml:"AddResponse"`
		}
	}

	if err := xml.Unmarshal(respData, &responseEnvelope); err != nil {
		return 0, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return responseEnvelope.Body.Response.AddResult, nil
}

func main() {
	// Create client
	client := NewCalculatorClient()

	// Test addition
	result, err := client.Add(5, 3)
	if err != nil {
		log.Fatalf("Failed to add numbers: %v", err)
	}

	fmt.Printf("5 + 3 = %d\n", result)

	// More examples
	result2, err := client.Add(100, 250)
	if err != nil {
		log.Fatalf("Failed to add numbers: %v", err)
	}

	fmt.Printf("100 + 250 = %d\n", result2)
}
