package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/thdev01/wsdl2api/internal/models"
)

// generateClientWithSecurity generates a SOAP client with WS-Security support
func (g *Generator) generateClientWithSecurity(def *models.Definitions) error {
	endpoint := g.findServiceEndpoint(def)

	content := fmt.Sprintf(`package %s

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/thdev01/wsdl2api/pkg/security"
)

// Client represents a SOAP client with WS-Security support
type Client struct {
	URL        string
	HTTPClient *http.Client
	Headers    map[string]string
	Security   *security.WSSecurity
	SOAPVersion string // "1.1" or "1.2"
}

// NewClient creates a new SOAP client
func NewClient(url string) *Client {
	if url == "" {
		url = "%s"
	}
	return &Client{
		URL:         url,
		HTTPClient:  &http.Client{},
		Headers:     make(map[string]string),
		SOAPVersion: "1.1",
	}
}

// SetBasicAuth sets basic authentication (WS-Security UsernameToken)
func (c *Client) SetBasicAuth(username, password string) {
	c.Security = &security.WSSecurity{
		Username:  username,
		Password:  password,
		UseDigest: false,
	}
}

// SetDigestAuth sets digest authentication (WS-Security UsernameToken with digest)
func (c *Client) SetDigestAuth(username, password string) {
	c.Security = &security.WSSecurity{
		Username:  username,
		Password:  password,
		UseDigest: true,
	}
}

// SetSOAPVersion sets the SOAP version (1.1 or 1.2)
func (c *Client) SetSOAPVersion(version string) {
	c.SOAPVersion = version
}

// SetHeader sets a custom HTTP header
func (c *Client) SetHeader(key, value string) {
	c.Headers[key] = value
}

// Call makes a SOAP call
func (c *Client) Call(soapAction string, request, response interface{}) error {
	// Build SOAP envelope based on version
	var envelope interface{}
	var contentType string

	if c.SOAPVersion == "1.2" {
		envelope = c.buildSOAP12Envelope(request)
		contentType = "application/soap+xml; charset=utf-8"
	} else {
		envelope = c.buildSOAP11Envelope(request)
		contentType = "text/xml; charset=utf-8"
	}

	// Marshal to XML
	xmlData, err := xml.MarshalIndent(envelope, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal request: %%w", err)
	}

	// Add XML header
	requestBody := []byte(xml.Header + string(xmlData))

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", c.URL, bytes.NewReader(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %%w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", contentType)
	if c.SOAPVersion == "1.1" {
		httpReq.Header.Set("SOAPAction", fmt.Sprintf("\"%%s\"", soapAction))
	}
	for key, value := range c.Headers {
		httpReq.Header.Set(key, value)
	}

	// Execute request
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to execute request: %%w", err)
	}
	defer resp.Body.Close()

	// Read response
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %%w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("SOAP request failed with status %%d: %%s", resp.StatusCode, string(respData))
	}

	// Parse SOAP response
	var responseEnvelope SOAPEnvelope
	responseEnvelope.Body.Content = response

	if err := xml.Unmarshal(respData, &responseEnvelope); err != nil {
		// Try SOAP 1.2 format
		var responseEnvelope12 SOAP12Envelope
		responseEnvelope12.Body.Content = response
		if err := xml.Unmarshal(respData, &responseEnvelope12); err != nil {
			return fmt.Errorf("failed to unmarshal response: %%w", err)
		}
	}

	return nil
}

// buildSOAP11Envelope builds a SOAP 1.1 envelope
func (c *Client) buildSOAP11Envelope(request interface{}) *SOAPEnvelope {
	envelope := &SOAPEnvelope{
		EnvNamespace: "http://schemas.xmlsoap.org/soap/envelope/",
		Body: SOAPBody{
			Content: request,
		},
	}

	// Add WS-Security header if configured
	if c.Security != nil {
		envelope.Header = &SOAPHeader{
			Security: security.NewSecurityHeader(c.Security),
		}
	}

	return envelope
}

// buildSOAP12Envelope builds a SOAP 1.2 envelope
func (c *Client) buildSOAP12Envelope(request interface{}) *SOAP12Envelope {
	envelope := &SOAP12Envelope{
		EnvNamespace: "http://www.w3.org/2003/05/soap-envelope",
		Body: SOAP12Body{
			Content: request,
		},
	}

	// Add WS-Security header if configured
	if c.Security != nil {
		envelope.Header = &SOAP12Header{
			Security: security.NewSecurityHeader(c.Security),
		}
	}

	return envelope
}

// SOAP 1.1 structures
type SOAPEnvelope struct {
	XMLName      xml.Name    ` + "`xml:\"soap:Envelope\"`" + `
	EnvNamespace string      ` + "`xml:\"xmlns:soap,attr\"`" + `
	Header       *SOAPHeader ` + "`xml:\"soap:Header,omitempty\"`" + `
	Body         SOAPBody    ` + "`xml:\"soap:Body\"`" + `
}

type SOAPHeader struct {
	XMLName  xml.Name                ` + "`xml:\"soap:Header\"`" + `
	Security *security.SecurityHeader ` + "`xml:\",omitempty\"`" + `
}

type SOAPBody struct {
	XMLName xml.Name    ` + "`xml:\"soap:Body\"`" + `
	Content interface{} ` + "`xml:\",innerxml\"`" + `
}

// SOAP 1.2 structures
type SOAP12Envelope struct {
	XMLName      xml.Name      ` + "`xml:\"env:Envelope\"`" + `
	EnvNamespace string        ` + "`xml:\"xmlns:env,attr\"`" + `
	Header       *SOAP12Header ` + "`xml:\"env:Header,omitempty\"`" + `
	Body         SOAP12Body    ` + "`xml:\"env:Body\"`" + `
}

type SOAP12Header struct {
	XMLName  xml.Name                ` + "`xml:\"env:Header\"`" + `
	Security *security.SecurityHeader ` + "`xml:\",omitempty\"`" + `
}

type SOAP12Body struct {
	XMLName xml.Name    ` + "`xml:\"env:Body\"`" + `
	Content interface{} ` + "`xml:\",innerxml\"`" + `
}

// SOAPFault represents a SOAP fault
type SOAPFault struct {
	XMLName xml.Name ` + "`xml:\"Fault\"`" + `
	Code    string   ` + "`xml:\"faultcode\"`" + `
	String  string   ` + "`xml:\"faultstring\"`" + `
	Actor   string   ` + "`xml:\"faultactor\"`" + `
	Detail  string   ` + "`xml:\"detail\"`" + `
}
`, g.packageName, endpoint)

	return os.WriteFile(filepath.Join(g.outputDir, "client.go"), []byte(content), 0644)
}
