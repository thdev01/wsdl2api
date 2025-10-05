package calculator

import "encoding/xml"

// Auto-generated types from WSDL

// AddRequest represents the request for Add operation
type AddRequest struct {
	XMLName    xml.Name `xml:"http://tempuri.org/ Add"`
	Parameters string   `xml:"parameters"`
}

// AddResponse represents the response for Add operation
type AddResponse struct {
	XMLName    xml.Name `xml:"http://tempuri.org/ AddResponse"`
	Parameters string   `xml:"parameters"`
}

// AddSoapInRequest for SOAP envelope
type AddSoapInRequest struct {
	XMLName    xml.Name `xml:"Add"`
	Parameters string   `xml:"parameters"`
}

// AddResult represents operation result
type AddResult struct {
	Result string
}
