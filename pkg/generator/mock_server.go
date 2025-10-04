package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/thdev01/wsdl2api/internal/models"
)

// generateMockServer generates a mock SOAP server for testing
func (g *Generator) generateMockServer(def *models.Definitions) error {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("package %s\n\n", g.packageName))
	b.WriteString(`import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// MockServer represents a mock SOAP server for testing
type MockServer struct {
	Port     int
	handlers map[string]MockHandler
}

// MockHandler is a function that handles a SOAP operation
type MockHandler func(request interface{}) (interface{}, error)

// NewMockServer creates a new mock server
func NewMockServer(port int) *MockServer {
	return &MockServer{
		Port:     port,
		handlers: make(map[string]MockHandler),
	}
}

// RegisterHandler registers a mock handler for an operation
func (m *MockServer) RegisterHandler(operation string, handler MockHandler) {
	m.handlers[operation] = handler
}

// Start starts the mock server
func (m *MockServer) Start() error {
	http.HandleFunc("/", m.handleSOAPRequest)

	addr := fmt.Sprintf(":%d", m.Port)
	log.Printf("Mock SOAP server listening on %s", addr)
	return http.ListenAndServe(addr, nil)
}

// handleSOAPRequest handles incoming SOAP requests
func (m *MockServer) handleSOAPRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request", http.StatusBadRequest)
		return
	}

	// Parse SOAP envelope to get operation name
	var envelope struct {
		XMLName xml.Name
		Body    struct {
			XMLName xml.Name
			Content string ` + "`xml:\",innerxml\"`" + `
		} ` + "`xml:\"Body\"`" + `
	}

	if err := xml.Unmarshal(body, &envelope); err != nil {
		m.sendSOAPFault(w, "Client", "Invalid SOAP envelope", "")
		return
	}

	// Extract operation name from body content
	operation := m.extractOperation(envelope.Body.Content)
	if operation == "" {
		m.sendSOAPFault(w, "Client", "Could not determine operation", "")
		return
	}

	// Find and execute handler
	handler, exists := m.handlers[operation]
	if !exists {
		m.sendSOAPFault(w, "Server", fmt.Sprintf("No mock handler for operation: %s", operation), "")
		return
	}

	// Execute mock handler (simplified - real implementation would unmarshal request)
	response, err := handler(nil)
	if err != nil {
		m.sendSOAPFault(w, "Server", err.Error(), "")
		return
	}

	// Send response
	m.sendSOAPResponse(w, response)
}

// extractOperation extracts the operation name from SOAP body content
func (m *MockServer) extractOperation(content string) string {
	// Simple XML parsing to get first element name
	content = strings.TrimSpace(content)
	if !strings.HasPrefix(content, "<") {
		return ""
	}

	end := strings.Index(content[1:], " ")
	if end == -1 {
		end = strings.Index(content[1:], ">")
	}

	if end == -1 {
		return ""
	}

	operation := content[1 : end+1]
	// Remove namespace prefix
	if idx := strings.Index(operation, ":"); idx != -1 {
		operation = operation[idx+1:]
	}

	return operation
}

// sendSOAPResponse sends a SOAP response
func (m *MockServer) sendSOAPResponse(w http.ResponseWriter, response interface{}) {
	envelope := struct {
		XMLName xml.Name    ` + "`xml:\"soap:Envelope\"`" + `
		Soap    string      ` + "`xml:\"xmlns:soap,attr\"`" + `
		Body    interface{} ` + "`xml:\"soap:Body\"`" + `
	}{
		Soap: "http://schemas.xmlsoap.org/soap/envelope/",
		Body: response,
	}

	xmlData, err := xml.MarshalIndent(envelope, "", "  ")
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(xml.Header))
	w.Write(xmlData)
}

// sendSOAPFault sends a SOAP fault
func (m *MockServer) sendSOAPFault(w http.ResponseWriter, code, message, detail string) {
	fault := struct {
		XMLName xml.Name ` + "`xml:\"soap:Envelope\"`" + `
		Soap    string   ` + "`xml:\"xmlns:soap,attr\"`" + `
		Body    struct {
			XMLName xml.Name ` + "`xml:\"soap:Body\"`" + `
			Fault   struct {
				XMLName     xml.Name ` + "`xml:\"soap:Fault\"`" + `
				Faultcode   string   ` + "`xml:\"faultcode\"`" + `
				Faultstring string   ` + "`xml:\"faultstring\"`" + `
				Detail      string   ` + "`xml:\"detail,omitempty\"`" + `
			}
		}
	}{
		Soap: "http://schemas.xmlsoap.org/soap/envelope/",
	}

	fault.Body.Fault.Faultcode = "soap:" + code
	fault.Body.Fault.Faultstring = message
	fault.Body.Fault.Detail = detail

	xmlData, _ := xml.MarshalIndent(fault, "", "  ")

	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(xml.Header))
	w.Write(xmlData)
}
`)

	// Generate default mock handlers for each operation
	b.WriteString("\n// Default mock handlers\n\n")

	for _, portType := range def.PortTypes {
		for _, op := range portType.Operations {
			methodName := toPascalCase(op.Name)
			outputMsg := g.findMessage(def, op.Output.Name)

			b.WriteString(fmt.Sprintf("// Mock%s is a default mock handler for %s operation\n", methodName, op.Name))
			b.WriteString(fmt.Sprintf("func Mock%s(request interface{}) (interface{}, error) {\n", methodName))
			b.WriteString("\t// TODO: Implement mock logic\n")
			b.WriteString(fmt.Sprintf("\treturn &%sResponse{}, nil\n", methodName))
			b.WriteString("}\n\n")
		}
	}

	// Generate example usage
	b.WriteString("// Example usage:\n/*\n")
	b.WriteString("func ExampleMockServer() {\n")
	b.WriteString("\tmock := NewMockServer(8080)\n\n")

	for _, portType := range def.PortTypes {
		if len(portType.Operations) > 0 {
			op := portType.Operations[0]
			methodName := toPascalCase(op.Name)
			b.WriteString(fmt.Sprintf("\t// Register custom handler for %s\n", op.Name))
			b.WriteString(fmt.Sprintf("\tmock.RegisterHandler(\"%s\", Mock%s)\n", op.Name, methodName))
			break
		}
	}

	b.WriteString("\n\tlog.Fatal(mock.Start())\n")
	b.WriteString("}\n*/\n")

	return os.WriteFile(filepath.Join(g.outputDir, "mock_server.go"), []byte(b.String()), 0644)
}
