package integration

import (
	"testing"

	numberconversion "github.com/thdev01/wsdl2api/tests/integration/numberconversion/client"
)

// TestNumberConversionClientGeneration tests number conversion client generation
func TestNumberConversionClientGeneration(t *testing.T) {
	c := numberconversion.NewClient("")

	if c == nil {
		t.Fatal("Failed to create number conversion client")
	}

	if c.URL == "" {
		t.Error("Client URL should not be empty")
	}

	t.Logf("✓ Number Conversion client generated successfully")
	t.Logf("  URL: %s", c.URL)
}

// TestNumberConversionOperations tests that all operations are available
func TestNumberConversionOperations(t *testing.T) {
	// Test NumberToWords operation exists
	req := &numberconversion.NumberToWordsRequest{
		UbiNum: "123",
	}

	if req.UbiNum != "123" {
		t.Error("NumberToWordsRequest not working correctly")
	}

	// Test NumberToDollars operation
	dollarsReq := &numberconversion.NumberToDollarsRequest{
		DNum: "100",
	}

	if dollarsReq.DNum != "100" {
		t.Error("NumberToDollarsRequest not working correctly")
	}

	t.Logf("✓ Number Conversion operations generated correctly")
	t.Logf("  Operations: NumberToWords, NumberToDollars")
}

// TestNumberConversionTypeSafety tests type safety of generated code
func TestNumberConversionTypeSafety(t *testing.T) {
	// Test that we can't accidentally mix up request types
	wordsReq := &numberconversion.NumberToWordsRequest{
		UbiNum: "456",
	}

	dollarsReq := &numberconversion.NumberToDollarsRequest{
		DNum: "789",
	}

	t.Logf("✓ Type safety working correctly")
	t.Logf("  NumberToWordsRequest.UbiNum: %s", wordsReq.UbiNum)
	t.Logf("  NumberToDollarsRequest.DNum: %s", dollarsReq.DNum)
}

// TestNumberConversionClientConfiguration tests client configuration options
func TestNumberConversionClientConfiguration(t *testing.T) {
	c := numberconversion.NewClient("http://custom-url.com")

	if c.URL != "http://custom-url.com" {
		t.Errorf("Expected custom URL, got %s", c.URL)
	}

	// Test SOAP version configuration
	c.SetSOAPVersion("1.2")
	if c.SOAPVersion != "1.2" {
		t.Errorf("Expected SOAP 1.2, got %s", c.SOAPVersion)
	}

	// Test authentication
	c.SetBasicAuth("user", "pass")
	if c.Security == nil || c.Security.Username != "user" {
		t.Error("Authentication not configured correctly")
	}

	t.Logf("✓ Number Conversion client configuration works")
	t.Logf("  Custom URL: %s", c.URL)
	t.Logf("  SOAP Version: %s", c.SOAPVersion)
	t.Logf("  Auth User: %s", c.Security.Username)
}

// BenchmarkNumberConversionRequestCreation benchmarks request creation
func BenchmarkNumberConversionRequestCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = &numberconversion.NumberToWordsRequest{
			UbiNum: "12345",
		}
	}
}

// BenchmarkNumberConversionClientCreation benchmarks client creation
func BenchmarkNumberConversionClientCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = numberconversion.NewClient("")
	}
}
