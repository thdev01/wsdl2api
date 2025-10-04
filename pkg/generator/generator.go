package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/thdev01/wsdl2api/internal/models"
)

// Generator generates Go code from WSDL definitions
type Generator struct {
	outputDir   string
	packageName string
}

// NewGenerator creates a new code generator
func NewGenerator(outputDir, packageName string) *Generator {
	return &Generator{
		outputDir:   outputDir,
		packageName: packageName,
	}
}

// Generate generates all code from WSDL definitions
func (g *Generator) Generate(def *models.Definitions) error {
	// Create output directory
	if err := os.MkdirAll(g.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate client with WS-Security support
	if err := g.generateClientWithSecurity(def); err != nil {
		return fmt.Errorf("failed to generate client: %w", err)
	}

	// Generate improved types
	if err := g.generateTypesImproved(def); err != nil {
		return fmt.Errorf("failed to generate types: %w", err)
	}

	// Generate operator functions
	if err := g.generateOperatorsImproved(def); err != nil {
		return fmt.Errorf("failed to generate operators: %w", err)
	}

	// Generate usage example
	if err := g.generateUsageExample(def); err != nil {
		return fmt.Errorf("failed to generate usage example: %w", err)
	}

	return nil
}

// GenerateWithMock generates all code including mock server
func (g *Generator) GenerateWithMock(def *models.Definitions) error {
	// Generate all standard code
	if err := g.Generate(def); err != nil {
		return err
	}

	// Generate mock server
	if err := g.generateMockServer(def); err != nil {
		return fmt.Errorf("failed to generate mock server: %w", err)
	}

	return nil
}

// generateClient generates the SOAP client code
func (g *Generator) generateClient(def *models.Definitions) error {
	tmpl := `package {{.PackageName}}

import (
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
}

// NewClient creates a new SOAP client
func NewClient(url string) *Client {
	return &Client{
		URL:        url,
		HTTPClient: &http.Client{},
	}
}

// Call makes a SOAP call
func (c *Client) Call(soapAction string, request interface{}, response interface{}) error {
	envelope := SOAPEnvelope{
		Body: SOAPBody{
			Content: request,
		},
	}

	// Marshal request
	reqData, err := xml.Marshal(envelope)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", c.URL, bytes.NewReader(reqData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "text/xml; charset=utf-8")
	httpReq.Header.Set("SOAPAction", soapAction)

	// Make request
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Unmarshal response
	var respEnvelope SOAPEnvelope
	if err := xml.Unmarshal(respData, &respEnvelope); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return xml.Unmarshal([]byte(respEnvelope.Body.Content), response)
}

// SOAP Envelope structures
type SOAPEnvelope struct {
	XMLName xml.Name ` + "`xml:\"http://schemas.xmlsoap.org/soap/envelope/ Envelope\"`" + `
	Body    SOAPBody
}

type SOAPBody struct {
	XMLName xml.Name ` + "`xml:\"http://schemas.xmlsoap.org/soap/envelope/ Body\"`" + `
	Content interface{}
}
`

	t, err := template.New("client").Parse(tmpl)
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(g.outputDir, "client.go"))
	if err != nil {
		return err
	}
	defer f.Close()

	return t.Execute(f, map[string]interface{}{
		"PackageName": g.packageName,
	})
}

// generateTypes generates Go structs for WSDL types
func (g *Generator) generateTypes(def *models.Definitions) error {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("package %s\n\n", g.packageName))
	b.WriteString("// Auto-generated types from WSDL\n\n")

	// Generate message types
	for _, msg := range def.Messages {
		structName := toPascalCase(msg.Name)
		b.WriteString(fmt.Sprintf("// %s represents %s message\n", structName, msg.Name))
		b.WriteString(fmt.Sprintf("type %s struct {\n", structName))

		for _, part := range msg.Parts {
			fieldName := toPascalCase(part.Name)
			fieldType := mapXSDTypeToGo(part.Type)
			b.WriteString(fmt.Sprintf("\t%s %s `xml:\"%s\"`\n", fieldName, fieldType, part.Name))
		}

		b.WriteString("}\n\n")
	}

	// Write to file
	return os.WriteFile(filepath.Join(g.outputDir, "types.go"), []byte(b.String()), 0644)
}

// generateOperations generates operation methods
func (g *Generator) generateOperations(def *models.Definitions) error {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("package %s\n\n", g.packageName))
	b.WriteString("import \"fmt\"\n\n")
	b.WriteString("// Auto-generated operations from WSDL\n\n")

	// Generate operations from port types
	for _, portType := range def.PortTypes {
		for _, op := range portType.Operations {
			methodName := toPascalCase(op.Name)
			inputType := toPascalCase(op.Input.Name)
			outputType := toPascalCase(op.Output.Name)

			// Find SOAP action from bindings
			soapAction := g.findSoapAction(def, op.Name)

			b.WriteString(fmt.Sprintf("// %s executes %s operation\n", methodName, op.Name))
			if op.Documentation != "" {
				b.WriteString(fmt.Sprintf("// %s\n", op.Documentation))
			}
			b.WriteString(fmt.Sprintf("func (c *Client) %s(req *%s) (*%s, error) {\n", methodName, inputType, outputType))
			b.WriteString(fmt.Sprintf("\tvar resp %s\n", outputType))
			b.WriteString(fmt.Sprintf("\terr := c.Call(\"%s\", req, &resp)\n", soapAction))
			b.WriteString("\tif err != nil {\n")
			b.WriteString(fmt.Sprintf("\t\treturn nil, fmt.Errorf(\"failed to execute %s: %%w\", err)\n", op.Name))
			b.WriteString("\t}\n")
			b.WriteString("\treturn &resp, nil\n")
			b.WriteString("}\n\n")
		}
	}

	// Write to file
	return os.WriteFile(filepath.Join(g.outputDir, "operations.go"), []byte(b.String()), 0644)
}

// findSoapAction finds the SOAP action for an operation
func (g *Generator) findSoapAction(def *models.Definitions, opName string) string {
	for _, binding := range def.Bindings {
		for _, op := range binding.Operations {
			if op.Name == opName {
				return op.SoapAction
			}
		}
	}
	return ""
}

// Helper functions
func toPascalCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	// Remove namespace prefix if present
	if idx := strings.LastIndex(s, ":"); idx != -1 {
		s = s[idx+1:]
	}

	// Split by common separators
	words := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == '.' || r == ' '
	})

	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + word[1:]
		}
	}

	return strings.Join(words, "")
}

func mapXSDTypeToGo(xsdType string) string {
	// Remove namespace prefix
	if idx := strings.LastIndex(xsdType, ":"); idx != -1 {
		xsdType = xsdType[idx+1:]
	}

	typeMap := map[string]string{
		"string":    "string",
		"int":       "int",
		"integer":   "int",
		"long":      "int64",
		"short":     "int16",
		"byte":      "byte",
		"boolean":   "bool",
		"float":     "float32",
		"double":    "float64",
		"decimal":   "float64",
		"dateTime":  "string",
		"date":      "string",
		"time":      "string",
		"base64Binary": "[]byte",
		"hexBinary": "[]byte",
	}

	if goType, ok := typeMap[xsdType]; ok {
		return goType
	}

	// If not a primitive type, assume it's a custom type
	return toPascalCase(xsdType)
}
