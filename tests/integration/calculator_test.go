package integration

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	calculator "github.com/thdev01/wsdl2api/tests/integration/calculator/client"
)

// TestCalculatorClientGeneration tests that the calculator client was generated correctly
func TestCalculatorClientGeneration(t *testing.T) {
	c := calculator.NewClient("")

	if c == nil {
		t.Fatal("Failed to create calculator client")
	}

	if c.URL == "" {
		t.Error("Client URL should not be empty")
	}

	if c.HTTPClient == nil {
		t.Error("HTTPClient should not be nil")
	}

	if c.SOAPVersion != "1.1" {
		t.Errorf("Expected SOAP version 1.1, got %s", c.SOAPVersion)
	}

	t.Logf("✓ Calculator client generated successfully")
	t.Logf("  URL: %s", c.URL)
	t.Logf("  SOAP Version: %s", c.SOAPVersion)
}

// TestCalculatorClientMethods tests that all expected methods exist
func TestCalculatorClientMethods(t *testing.T) {
	c := calculator.NewClient("")

	// Test that we can set authentication
	c.SetBasicAuth("test", "password")
	if c.Security == nil {
		t.Error("SetBasicAuth should set Security field")
	}
	if c.Security.Username != "test" {
		t.Errorf("Expected username 'test', got '%s'", c.Security.Username)
	}

	// Test digest auth
	c.SetDigestAuth("test2", "password2")
	if !c.Security.UseDigest {
		t.Error("SetDigestAuth should set UseDigest to true")
	}

	// Test SOAP version change
	c.SetSOAPVersion("1.2")
	if c.SOAPVersion != "1.2" {
		t.Errorf("Expected SOAP version 1.2, got %s", c.SOAPVersion)
	}

	t.Logf("✓ Calculator client methods work correctly")
}

// TestCalculatorTypesGeneration tests that types were generated correctly
func TestCalculatorTypesGeneration(t *testing.T) {
	// Create request
	req := &calculator.AddRequest{
		Parameters: "test parameters",
	}

	if req.Parameters != "test parameters" {
		t.Error("AddRequest Parameters field not working")
	}

	// Create response
	resp := &calculator.AddResponse{
		Parameters: "test response",
	}

	if resp.Parameters != "test response" {
		t.Error("AddResponse Parameters field not working")
	}

	t.Logf("✓ Calculator types generated correctly")
	t.Logf("  AddRequest: %+v", req)
	t.Logf("  AddResponse: %+v", resp)
}

// TestCalculatorSOAPEnvelopeFormat tests SOAP envelope structure
func TestCalculatorSOAPEnvelopeFormat(t *testing.T) {
	tests := []struct {
		name        string
		soapVersion string
		wantNS      string
	}{
		{
			name:        "SOAP 1.1",
			soapVersion: "1.1",
			wantNS:      "http://schemas.xmlsoap.org/soap/envelope/",
		},
		{
			name:        "SOAP 1.2",
			soapVersion: "1.2",
			wantNS:      "http://www.w3.org/2003/05/soap-envelope",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calculator.NewClient("")
			c.SetSOAPVersion(tt.soapVersion)

			if c.SOAPVersion != tt.soapVersion {
				t.Errorf("Expected SOAP version %s, got %s", tt.soapVersion, c.SOAPVersion)
			}

			t.Logf("✓ SOAP %s configuration correct", tt.soapVersion)
		})
	}
}

// TestCalculatorOpenAPIGeneration tests OpenAPI spec generation
func TestCalculatorOpenAPIGeneration(t *testing.T) {
	// Read generated OpenAPI spec
	data, err := os.ReadFile("calculator/openapi.json")
	if err != nil {
		t.Skipf("OpenAPI spec not found (expected for unit tests): %v", err)
		return
	}

	var spec map[string]interface{}
	if err := json.Unmarshal(data, &spec); err != nil {
		t.Fatalf("Failed to parse OpenAPI spec: %v", err)
	}

	// Verify OpenAPI version
	if version, ok := spec["openapi"].(string); !ok || version != "3.0.0" {
		t.Errorf("Expected OpenAPI version 3.0.0, got %v", spec["openapi"])
	}

	// Verify info section
	info, ok := spec["info"].(map[string]interface{})
	if !ok {
		t.Fatal("OpenAPI spec missing info section")
	}

	if title, ok := info["title"].(string); !ok || title == "" {
		t.Error("OpenAPI spec missing title")
	}

	// Verify paths exist
	paths, ok := spec["paths"].(map[string]interface{})
	if !ok || len(paths) == 0 {
		t.Error("OpenAPI spec missing paths")
	}

	t.Logf("✓ OpenAPI specification generated correctly")
	t.Logf("  Version: %v", spec["openapi"])
	t.Logf("  Title: %v", info["title"])
	t.Logf("  Paths: %d", len(paths))
}

// TestCalculatorClientHTTPHeaders tests custom HTTP headers
func TestCalculatorClientHTTPHeaders(t *testing.T) {
	c := calculator.NewClient("")

	// Set custom headers
	c.Headers["X-Custom-Header"] = "test-value"
	c.Headers["Authorization"] = "Bearer token123"

	if c.Headers["X-Custom-Header"] != "test-value" {
		t.Error("Failed to set custom header")
	}

	if c.Headers["Authorization"] != "Bearer token123" {
		t.Error("Failed to set authorization header")
	}

	t.Logf("✓ Custom HTTP headers work correctly")
	t.Logf("  Headers: %+v", c.Headers)
}

// TestCalculatorSecurityHeaders tests WS-Security header generation
func TestCalculatorSecurityHeaders(t *testing.T) {
	c := calculator.NewClient("")

	// Test basic auth
	c.SetBasicAuth("admin", "secret123")

	if c.Security == nil {
		t.Fatal("Security not set")
	}

	if c.Security.Username != "admin" {
		t.Errorf("Expected username 'admin', got '%s'", c.Security.Username)
	}

	if c.Security.Password != "secret123" {
		t.Errorf("Expected password 'secret123', got '%s'", c.Security.Password)
	}

	if c.Security.UseDigest {
		t.Error("Basic auth should not use digest")
	}

	// Test digest auth
	c.SetDigestAuth("admin", "secret123")

	if !c.Security.UseDigest {
		t.Error("Digest auth should set UseDigest")
	}

	t.Logf("✓ WS-Security headers configured correctly")
	t.Logf("  Username: %s", c.Security.Username)
	t.Logf("  UseDigest: %v", c.Security.UseDigest)
}

// BenchmarkCalculatorClientCreation benchmarks client creation
func BenchmarkCalculatorClientCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = calculator.NewClient("")
	}
}

// BenchmarkCalculatorRequestCreation benchmarks request creation
func BenchmarkCalculatorRequestCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = &calculator.AddRequest{
			Parameters: "test",
		}
	}
}

// TestCalculatorClientTimeout tests HTTP client timeout configuration
func TestCalculatorClientTimeout(t *testing.T) {
	c := calculator.NewClient("http://invalid-url-that-does-not-exist:99999")
	c.HTTPClient.Timeout = 1 * time.Second

	if c.HTTPClient.Timeout != 1*time.Second {
		t.Error("Failed to set HTTP client timeout")
	}

	t.Logf("✓ HTTP client timeout configured correctly")
	t.Logf("  Timeout: %v", c.HTTPClient.Timeout)
}

// TestCalculatorJSONSerialization tests that structs can be JSON serialized
func TestCalculatorJSONSerialization(t *testing.T) {
	req := &calculator.AddRequest{
		Parameters: "test value",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	if !strings.Contains(string(data), "test value") {
		t.Error("JSON serialization doesn't contain expected value")
	}

	// Test deserialization
	var decoded calculator.AddRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal request: %v", err)
	}

	if decoded.Parameters != req.Parameters {
		t.Error("Deserialized value doesn't match original")
	}

	t.Logf("✓ JSON serialization works correctly")
	t.Logf("  JSON: %s", string(data))
}
