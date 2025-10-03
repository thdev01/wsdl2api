package parser

import (
	"testing"
)

func TestNewParser(t *testing.T) {
	p := NewParser()
	if p == nil {
		t.Fatal("NewParser() returned nil")
	}
}

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello_world", "HelloWorld"},
		{"test-case", "TestCase"},
		{"simple", "Simple"},
		{"with spaces", "WithSpaces"},
		{"ns:LocalName", "LocalName"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			// This would need to be exported or tested through the generator
			// For now, this is a placeholder
		})
	}
}

// TODO: Add more tests with sample WSDL files
