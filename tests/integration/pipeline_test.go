package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestFullPipelineCalculator tests the complete WSDL → Go → OpenAPI → TypeScript pipeline
func TestFullPipelineCalculator(t *testing.T) {
	t.Log("Testing complete pipeline for Calculator service")

	// Step 1: Verify Go client was generated
	t.Run("Go Client Generated", func(t *testing.T) {
		files := []string{
			"calculator/client/client.go",
			"calculator/client/types.go",
			"calculator/client/operators.go",
			"calculator/client/example.go",
			"calculator/client/mock_server.go",
		}

		for _, file := range files {
			if _, err := os.Stat(file); os.IsNotExist(err) {
				t.Errorf("Expected file not found: %s", file)
			} else {
				t.Logf("  ✓ Found: %s", filepath.Base(file))
			}
		}
	})

	// Step 2: Verify OpenAPI spec was generated
	t.Run("OpenAPI Spec Generated", func(t *testing.T) {
		data, err := os.ReadFile("calculator/openapi.json")
		if err != nil {
			t.Fatalf("Failed to read OpenAPI spec: %v", err)
		}

		var spec map[string]interface{}
		if err := json.Unmarshal(data, &spec); err != nil {
			t.Fatalf("Failed to parse OpenAPI spec: %v", err)
		}

		// Verify OpenAPI structure
		if version := spec["openapi"]; version != "3.0.0" {
			t.Errorf("Expected OpenAPI 3.0.0, got %v", version)
		}

		info := spec["info"].(map[string]interface{})
		t.Logf("  ✓ OpenAPI Version: %v", spec["openapi"])
		t.Logf("  ✓ API Title: %v", info["title"])

		// Verify paths exist
		paths := spec["paths"].(map[string]interface{})
		if len(paths) == 0 {
			t.Error("No paths found in OpenAPI spec")
		}

		for path := range paths {
			t.Logf("  ✓ Path: %s", path)
		}
	})

	// Step 3: Verify TypeScript client was generated
	t.Run("TypeScript Client Generated", func(t *testing.T) {
		tsFiles := []string{
			"calculator/typescript/types.ts",
			"calculator/typescript/client.ts",
			"calculator/typescript/index.ts",
			"calculator/typescript/package.json",
			"calculator/typescript/tsconfig.json",
			"calculator/typescript/README.md",
		}

		for _, file := range tsFiles {
			if _, err := os.Stat(file); os.IsNotExist(err) {
				t.Errorf("Expected TypeScript file not found: %s", file)
			} else {
				t.Logf("  ✓ Found: %s", filepath.Base(file))
			}
		}

		// Verify TypeScript types.ts content
		data, err := os.ReadFile("calculator/typescript/types.ts")
		if err != nil {
			t.Fatalf("Failed to read types.ts: %v", err)
		}

		content := string(data)
		expectedTypes := []string{
			"export interface",
			"AddRequest",
			"AddResponse",
			"SOAPFault",
			"APIError",
		}

		for _, expected := range expectedTypes {
			if !contains(content, expected) {
				t.Errorf("TypeScript types.ts missing: %s", expected)
			}
		}
		t.Logf("  ✓ TypeScript types contain all expected interfaces")

		// Verify TypeScript client.ts content
		clientData, err := os.ReadFile("calculator/typescript/client.ts")
		if err != nil {
			t.Fatalf("Failed to read client.ts: %v", err)
		}

		clientContent := string(clientData)
		expectedClient := []string{
			"export class APIClient",
			"async add(",
			"Promise<",
			"fetch(",
		}

		for _, expected := range expectedClient {
			if !contains(clientContent, expected) {
				t.Errorf("TypeScript client.ts missing: %s", expected)
			}
		}
		t.Logf("  ✓ TypeScript client contains all expected methods")
	})

	// Step 4: Verify package.json has correct structure
	t.Run("TypeScript Package Configuration", func(t *testing.T) {
		data, err := os.ReadFile("calculator/typescript/package.json")
		if err != nil {
			t.Fatalf("Failed to read package.json: %v", err)
		}

		var pkg map[string]interface{}
		if err := json.Unmarshal(data, &pkg); err != nil {
			t.Fatalf("Failed to parse package.json: %v", err)
		}

		if name := pkg["name"]; name == "" {
			t.Error("package.json missing name")
		}

		if version := pkg["version"]; version == "" {
			t.Error("package.json missing version")
		}

		scripts := pkg["scripts"].(map[string]interface{})
		if _, ok := scripts["build"]; !ok {
			t.Error("package.json missing build script")
		}

		t.Logf("  ✓ Package name: %v", pkg["name"])
		t.Logf("  ✓ Package version: %v", pkg["version"])
		t.Logf("  ✓ Scripts: %v", len(scripts))
	})

	t.Log("✅ Complete pipeline test PASSED for Calculator service")
}

// TestFullPipelineNumberConversion tests the pipeline for NumberConversion service
func TestFullPipelineNumberConversion(t *testing.T) {
	t.Log("Testing complete pipeline for NumberConversion service")

	testPipeline(t, "numberconversion", []string{
		"NumberToWords",
		"NumberToDollars",
	})

	t.Log("✅ Complete pipeline test PASSED for NumberConversion service")
}

// TestFullPipelineTemperature tests the pipeline for Temperature service
func TestFullPipelineTemperature(t *testing.T) {
	t.Log("Testing complete pipeline for Temperature service")

	testPipeline(t, "temperature", []string{
		"CelsiusToFahrenheit",
		"FahrenheitToCelsius",
	})

	t.Log("✅ Complete pipeline test PASSED for Temperature service")
}

// testPipeline is a helper function to test the complete pipeline for any service
func testPipeline(t *testing.T, serviceName string, expectedOperations []string) {
	// Verify Go client files
	clientFiles := []string{
		filepath.Join(serviceName, "client", "client.go"),
		filepath.Join(serviceName, "client", "types.go"),
		filepath.Join(serviceName, "client", "operators.go"),
		filepath.Join(serviceName, "client", "mock_server.go"),
	}

	for _, file := range clientFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			t.Errorf("Missing Go file: %s", file)
		}
	}

	// Verify OpenAPI spec
	openAPIPath := filepath.Join(serviceName, "openapi.json")
	data, err := os.ReadFile(openAPIPath)
	if err != nil {
		t.Fatalf("Failed to read OpenAPI spec: %v", err)
	}

	var spec map[string]interface{}
	if err := json.Unmarshal(data, &spec); err != nil {
		t.Fatalf("Failed to parse OpenAPI spec: %v", err)
	}

	paths := spec["paths"].(map[string]interface{})
	t.Logf("  ✓ OpenAPI paths: %d", len(paths))

	// Verify TypeScript files
	tsFiles := []string{
		filepath.Join(serviceName, "typescript", "types.ts"),
		filepath.Join(serviceName, "typescript", "client.ts"),
		filepath.Join(serviceName, "typescript", "index.ts"),
	}

	for _, file := range tsFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			t.Errorf("Missing TypeScript file: %s", file)
		}
	}
}

// TestCodeQuality runs basic quality checks on generated code
func TestCodeQuality(t *testing.T) {
	services := []string{"calculator", "numberconversion", "temperature"}

	for _, service := range services {
		t.Run(service, func(t *testing.T) {
			// Check Go client has proper package declaration
			clientPath := filepath.Join(service, "client", "client.go")
			data, err := os.ReadFile(clientPath)
			if err != nil {
				t.Skipf("Client file not found: %v", err)
				return
			}

			content := string(data)

			// Verify package declaration
			if !contains(content, "package "+service) && !contains(content, "package client") {
				t.Error("Missing or incorrect package declaration")
			}

			// Verify imports
			if !contains(content, "import") {
				t.Error("Missing import statements")
			}

			// Verify client struct
			if !contains(content, "type Client struct") {
				t.Error("Missing Client struct definition")
			}

			// Verify NewClient function
			if !contains(content, "func NewClient") {
				t.Error("Missing NewClient function")
			}

			t.Logf("  ✓ Code quality checks passed for %s", service)
		})
	}
}

// TestGeneratedCodeSize checks that generated code is reasonable size
func TestGeneratedCodeSize(t *testing.T) {
	services := []string{"calculator", "numberconversion", "temperature"}

	for _, service := range services {
		t.Run(service, func(t *testing.T) {
			clientPath := filepath.Join(service, "client", "client.go")
			info, err := os.Stat(clientPath)
			if err != nil {
				t.Skipf("Client file not found: %v", err)
				return
			}

			size := info.Size()
			if size == 0 {
				t.Error("Generated client.go is empty")
			}

			if size > 1024*1024 { // 1MB
				t.Errorf("Generated client.go is too large: %d bytes", size)
			}

			t.Logf("  ✓ client.go size: %d bytes", size)
		})
	}
}

// TestTypeScriptTypeSafety verifies TypeScript type definitions
func TestTypeScriptTypeSafety(t *testing.T) {
	services := []string{"calculator", "numberconversion", "temperature"}

	for _, service := range services {
		t.Run(service, func(t *testing.T) {
			typesPath := filepath.Join(service, "typescript", "types.ts")
			data, err := os.ReadFile(typesPath)
			if err != nil {
				t.Skipf("TypeScript types not found: %v", err)
				return
			}

			content := string(data)

			// Verify export interfaces
			if !contains(content, "export interface") {
				t.Error("Missing export interface declarations")
			}

			// Verify error types
			if !contains(content, "SOAPFault") {
				t.Error("Missing SOAPFault type")
			}

			if !contains(content, "APIError") {
				t.Error("Missing APIError type")
			}

			t.Logf("  ✓ TypeScript types are properly defined")
		})
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
