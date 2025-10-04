package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/thdev01/wsdl2api/internal/models"
)

// generateClientImproved generates an improved SOAP client with proper XML handling
func (g *Generator) generateClientImproved(def *models.Definitions) error {
	var b strings.Builder

	// Find service endpoint
	endpoint := g.findServiceEndpoint(def)

	b.WriteString(fmt.Sprintf("package %s\n\n", g.packageName))
	b.WriteString(`import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

// Client represents a SOAP client
type Client struct {
	URL        string
	HTTPClient *http.Client
	Headers    map[string]string
}

// NewClient creates a new SOAP client
func NewClient(url string) *Client {
	if url == "" {
		url = "` + endpoint + `"
	}
	return &Client{
		URL:        url,
		HTTPClient: &http.Client{},
		Headers:    make(map[string]string),
	}
}

// SetHeader sets a custom HTTP header
func (c *Client) SetHeader(key, value string) {
	c.Headers[key] = value
}

// Call makes a SOAP call
func (c *Client) Call(soapAction string, request, response interface{}) error {
	// Build SOAP envelope
	envelope := &SOAPEnvelope{
		EnvNamespace: "http://schemas.xmlsoap.org/soap/envelope/",
		Body: SOAPBody{
			Content: request,
		},
	}

	// Marshal to XML
	xmlData, err := xml.MarshalIndent(envelope, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Add XML header
	requestBody := []byte(xml.Header + string(xmlData))

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", c.URL, bytes.NewReader(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "text/xml; charset=utf-8")
	httpReq.Header.Set("SOAPAction", fmt.Sprintf("\"%s\"", soapAction))
	for key, value := range c.Headers {
		httpReq.Header.Set(key, value)
	}

	// Execute request
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("SOAP request failed with status %d: %s", resp.StatusCode, string(respData))
	}

	// Parse SOAP response
	var responseEnvelope SOAPEnvelope
	responseEnvelope.Body.Content = response

	if err := xml.Unmarshal(respData, &responseEnvelope); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}

// SOAPEnvelope represents a SOAP 1.1 envelope
type SOAPEnvelope struct {
	XMLName      xml.Name ` + "`xml:\"soap:Envelope\"`" + `
	EnvNamespace string   ` + "`xml:\"xmlns:soap,attr\"`" + `
	Body         SOAPBody ` + "`xml:\"soap:Body\"`" + `
}

// SOAPBody represents the SOAP body
type SOAPBody struct {
	XMLName xml.Name    ` + "`xml:\"soap:Body\"`" + `
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
`)

	return os.WriteFile(filepath.Join(g.outputDir, "client.go"), []byte(b.String()), 0644)
}

// generateOperatorsImproved generates easy-to-use operator functions
func (g *Generator) generateOperatorsImproved(def *models.Definitions) error {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("package %s\n\n", g.packageName))
	b.WriteString("import \"fmt\"\n\n")
	b.WriteString("// Auto-generated operator functions for easy usage\n\n")

	// Find target namespace
	targetNS := def.TargetNamespace

	// Generate operators for each operation
	for _, portType := range def.PortTypes {
		for _, op := range portType.Operations {
			methodName := toPascalCase(op.Name)
			soapAction := g.findSoapAction(def, op.Name)

			// Find input/output message details
			inputMsg := g.findMessage(def, op.Input.Name)
			outputMsg := g.findMessage(def, op.Output.Name)

			if inputMsg == nil || outputMsg == nil {
				continue
			}

			// Generate parameter list
			params := g.generateParams(inputMsg)
			inputStruct := g.generateInputStruct(inputMsg, targetNS)
			outputField := g.generateOutputField(outputMsg)

			// Generate operator function
			b.WriteString(fmt.Sprintf("// %s is an easy-to-use operator for the %s operation\n", methodName, op.Name))
			if op.Documentation != "" {
				b.WriteString(fmt.Sprintf("// %s\n", op.Documentation))
			}
			b.WriteString(fmt.Sprintf("func (c *Client) %s(%s) (%s, error) {\n", methodName, params, outputField))
			b.WriteString(fmt.Sprintf("\trequest := %s\n", inputStruct))
			b.WriteString(fmt.Sprintf("\tvar response %sResponse\n\n", methodName))
			b.WriteString(fmt.Sprintf("\terr := c.Call(\"%s\", request, &response)\n", soapAction))
			b.WriteString("\tif err != nil {\n")
			b.WriteString(fmt.Sprintf("\t\treturn %s, fmt.Errorf(\"failed to execute %s: %%w\", err)\n", g.getZeroValue(outputField), op.Name))
			b.WriteString("\t}\n\n")
			b.WriteString(fmt.Sprintf("\treturn response.%sResult, nil\n", methodName))
			b.WriteString("}\n\n")
		}
	}

	return os.WriteFile(filepath.Join(g.outputDir, "operators.go"), []byte(b.String()), 0644)
}

// generateTypesImproved generates improved type definitions with proper XML tags
func (g *Generator) generateTypesImproved(def *models.Definitions) error {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("package %s\n\n", g.packageName))
	b.WriteString("import \"encoding/xml\"\n\n")
	b.WriteString("// Auto-generated types from WSDL\n\n")

	targetNS := def.TargetNamespace

	// Generate request/response types for each operation
	for _, portType := range def.PortTypes {
		for _, op := range portType.Operations {
			methodName := toPascalCase(op.Name)

			// Find messages
			inputMsg := g.findMessage(def, op.Input.Name)
			outputMsg := g.findMessage(def, op.Output.Name)

			if inputMsg == nil || outputMsg == nil {
				continue
			}

			// Generate request type
			b.WriteString(fmt.Sprintf("// %sRequest represents the request for %s operation\n", methodName, op.Name))
			b.WriteString(fmt.Sprintf("type %sRequest struct {\n", methodName))
			b.WriteString(fmt.Sprintf("\tXMLName xml.Name `xml:\"%s %s\"`\n", targetNS, op.Name))

			for _, part := range inputMsg.Parts {
				fieldName := toPascalCase(part.Name)
				fieldType := mapXSDTypeToGo(part.Type)
				xmlTag := part.Name
				b.WriteString(fmt.Sprintf("\t%s %s `xml:\"%s\"`\n", fieldName, fieldType, xmlTag))
			}
			b.WriteString("}\n\n")

			// Generate response type
			b.WriteString(fmt.Sprintf("// %sResponse represents the response for %s operation\n", methodName, op.Name))
			b.WriteString(fmt.Sprintf("type %sResponse struct {\n", methodName))
			b.WriteString(fmt.Sprintf("\tXMLName xml.Name `xml:\"%s %sResponse\"`\n", targetNS, op.Name))

			for _, part := range outputMsg.Parts {
				fieldName := toPascalCase(part.Name)
				fieldType := mapXSDTypeToGo(part.Type)
				xmlTag := part.Name
				b.WriteString(fmt.Sprintf("\t%s %s `xml:\"%s\"`\n", fieldName, fieldType, xmlTag))
			}
			b.WriteString("}\n\n")
		}
	}

	return os.WriteFile(filepath.Join(g.outputDir, "types.go"), []byte(b.String()), 0644)
}

// Helper methods

func (g *Generator) findMessage(def *models.Definitions, name string) *models.Message {
	// Remove namespace prefix
	if idx := strings.LastIndex(name, ":"); idx != -1 {
		name = name[idx+1:]
	}

	for _, msg := range def.Messages {
		msgName := msg.Name
		if idx := strings.LastIndex(msgName, ":"); idx != -1 {
			msgName = msgName[idx+1:]
		}
		if msgName == name {
			return &msg
		}
	}
	return nil
}

func (g *Generator) findServiceEndpoint(def *models.Definitions) string {
	for _, svc := range def.Services {
		for _, port := range svc.Ports {
			if port.Address != "" {
				return port.Address
			}
		}
	}
	return "http://localhost:8080/service"
}

func (g *Generator) generateParams(msg *models.Message) string {
	var params []string
	for _, part := range msg.Parts {
		fieldName := strings.ToLower(string(part.Name[0])) + part.Name[1:]
		fieldType := mapXSDTypeToGo(part.Type)
		params = append(params, fmt.Sprintf("%s %s", fieldName, fieldType))
	}
	return strings.Join(params, ", ")
}

func (g *Generator) generateInputStruct(msg *models.Message, targetNS string) string {
	var fields []string
	for _, part := range msg.Parts {
		fieldName := toPascalCase(part.Name)
		value := strings.ToLower(string(part.Name[0])) + part.Name[1:]
		fields = append(fields, fmt.Sprintf("%s: %s", fieldName, value))
	}
	return fmt.Sprintf("&%sRequest{%s}", toPascalCase(msg.Name), strings.Join(fields, ", "))
}

func (g *Generator) generateOutputField(msg *models.Message) string {
	if len(msg.Parts) > 0 {
		return mapXSDTypeToGo(msg.Parts[0].Type)
	}
	return "interface{}"
}

func (g *Generator) getZeroValue(typeName string) string {
	switch typeName {
	case "string":
		return "\"\""
	case "int", "int64", "int32", "int16":
		return "0"
	case "float32", "float64":
		return "0.0"
	case "bool":
		return "false"
	default:
		return "nil"
	}
}
