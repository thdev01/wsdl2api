package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/thdev01/wsdl2api/internal/models"
)

// generateUsageExample generates an example file showing how to use the client
func (g *Generator) generateUsageExample(def *models.Definitions) error {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("package %s\n\n", g.packageName))
	b.WriteString(`// This file contains usage examples for the generated SOAP client
// To use this client in your code:
//
// import "your-module/` + g.outputDir + `"
//
// Example usage:

/*
package main

import (
	"fmt"
	"log"

	"` + g.outputDir + `"
)

func main() {
	// Create a new client
	client := ` + g.packageName + `.NewClient("")

	// You can also specify a custom URL:
	// client := ` + g.packageName + `.NewClient("http://your-service-url")

`)

	// Generate example for first operation
	if len(def.PortTypes) > 0 && len(def.PortTypes[0].Operations) > 0 {
		op := def.PortTypes[0].Operations[0]
		methodName := toPascalCase(op.Name)
		inputMsg := g.findMessage(def, op.Input.Name)

		if inputMsg != nil && len(inputMsg.Parts) > 0 {
			// Generate example parameters
			var exampleParams []string
			for _, part := range inputMsg.Parts {
				exampleValue := g.getExampleValue(mapXSDTypeToGo(part.Type))
				exampleParams = append(exampleParams, exampleValue)
			}

			b.WriteString(fmt.Sprintf("\t// Example: Call %s operation\n", op.Name))
			b.WriteString(fmt.Sprintf("\tresult, err := client.%s(%s)\n", methodName, strings.Join(exampleParams, ", ")))
			b.WriteString("\tif err != nil {\n")
			b.WriteString(fmt.Sprintf("\t\tlog.Fatalf(\"Failed to call %s: %%v\", err)\n", op.Name))
			b.WriteString("\t}\n\n")
			b.WriteString("\tfmt.Printf(\"Result: %+v\\n\", result)\n")
		}
	}

	b.WriteString("}\n*/\n\n")

	// Add quick reference for all operations
	b.WriteString("// Available Operations:\n//\n")
	for _, portType := range def.PortTypes {
		for _, op := range portType.Operations {
			methodName := toPascalCase(op.Name)
			inputMsg := g.findMessage(def, op.Input.Name)

			if inputMsg != nil {
				params := g.generateParams(inputMsg)
				outputMsg := g.findMessage(def, op.Output.Name)
				outputType := "interface{}"
				if outputMsg != nil && len(outputMsg.Parts) > 0 {
					outputType = mapXSDTypeToGo(outputMsg.Parts[0].Type)
				}

				b.WriteString(fmt.Sprintf("// client.%s(%s) (%s, error)\n", methodName, params, outputType))
				if op.Documentation != "" {
					b.WriteString(fmt.Sprintf("//   %s\n", op.Documentation))
				}
				b.WriteString("//\n")
			}
		}
	}

	return os.WriteFile(filepath.Join(g.outputDir, "example.go"), []byte(b.String()), 0644)
}

// getExampleValue returns an example value for a Go type
func (g *Generator) getExampleValue(goType string) string {
	switch goType {
	case "string":
		return "\"example\""
	case "int", "int64", "int32", "int16":
		return "42"
	case "float32", "float64":
		return "3.14"
	case "bool":
		return "true"
	default:
		return "nil"
	}
}
