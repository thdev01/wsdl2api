package numberconversion

import "encoding/xml"

// Auto-generated types from WSDL

// NumberToWordsRequest represents NumberToWords request
type NumberToWordsRequest struct {
	XMLName xml.Name `xml:"NumberToWords"`
	UbiNum  string   `xml:"ubiNum"`
}

// NumberToWordsResponse represents NumberToWords response
type NumberToWordsResponse struct {
	XMLName               xml.Name `xml:"NumberToWordsResponse"`
	NumberToWordsResult   string   `xml:"NumberToWordsResult"`
}

// NumberToDollarsRequest represents NumberToDollars request
type NumberToDollarsRequest struct {
	XMLName xml.Name `xml:"NumberToDollars"`
	DNum    string   `xml:"dNum"`
}

// NumberToDollarsResponse represents NumberToDollars response
type NumberToDollarsResponse struct {
	XMLName                 xml.Name `xml:"NumberToDollarsResponse"`
	NumberToDollarsResult   string   `xml:"NumberToDollarsResult"`
}
