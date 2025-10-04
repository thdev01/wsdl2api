package typescript

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/thdev01/wsdl2api/pkg/exporter"
)

// Generator generates TypeScript client code from OpenAPI spec
type Generator struct {
	outputDir string
	spec      *exporter.OpenAPISpec
}

// NewGenerator creates a new TypeScript generator
func NewGenerator(outputDir string, spec *exporter.OpenAPISpec) *Generator {
	return &Generator{
		outputDir: outputDir,
		spec:      spec,
	}
}

// Generate generates TypeScript client code
func (g *Generator) Generate() error {
	// Create output directory
	if err := os.MkdirAll(g.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate types
	if err := g.generateTypes(); err != nil {
		return fmt.Errorf("failed to generate types: %w", err)
	}

	// Generate API client
	if err := g.generateClient(); err != nil {
		return fmt.Errorf("failed to generate client: %w", err)
	}

	// Generate index file
	if err := g.generateIndex(); err != nil {
		return fmt.Errorf("failed to generate index: %w", err)
	}

	// Generate package.json
	if err := g.generatePackageJSON(); err != nil {
		return fmt.Errorf("failed to generate package.json: %w", err)
	}

	// Generate tsconfig.json
	if err := g.generateTSConfig(); err != nil {
		return fmt.Errorf("failed to generate tsconfig.json: %w", err)
	}

	// Generate README
	if err := g.generateReadme(); err != nil {
		return fmt.Errorf("failed to generate README: %w", err)
	}

	return nil
}

// generateTypes generates TypeScript types from OpenAPI schemas
func (g *Generator) generateTypes() error {
	var b strings.Builder

	b.WriteString("// Auto-generated TypeScript types from OpenAPI specification\n\n")

	// Generate request/response types from paths
	for _, pathItem := range g.spec.Paths {
		if pathItem.Post != nil {
			op := pathItem.Post

			// Generate request type
			if op.RequestBody != nil {
				typeName := toPascalCase(op.OperationID) + "Request"
				b.WriteString(g.generateTypeFromSchema(typeName, op.RequestBody.Content["application/json"].Schema))
			}

			// Generate response type
			if resp, ok := op.Responses["200"]; ok {
				if content, ok := resp.Content["application/json"]; ok {
					typeName := toPascalCase(op.OperationID) + "Response"
					b.WriteString(g.generateTypeFromSchema(typeName, content.Schema))
				}
			}
		}
	}

	// Generate error types
	b.WriteString("// Error types\n")
	b.WriteString("export interface SOAPFault {\n")
	b.WriteString("  faultcode: string;\n")
	b.WriteString("  faultstring: string;\n")
	b.WriteString("  detail?: string;\n")
	b.WriteString("}\n\n")

	b.WriteString("export interface APIError {\n")
	b.WriteString("  message: string;\n")
	b.WriteString("  status: number;\n")
	b.WriteString("  fault?: SOAPFault;\n")
	b.WriteString("}\n\n")

	return os.WriteFile(filepath.Join(g.outputDir, "types.ts"), []byte(b.String()), 0644)
}

// generateTypeFromSchema generates a TypeScript type from OpenAPI schema
func (g *Generator) generateTypeFromSchema(name string, schema *exporter.OpenAPISchema) string {
	var b strings.Builder

	if schema == nil {
		return ""
	}

	b.WriteString(fmt.Sprintf("export interface %s {\n", name))

	if schema.Properties != nil {
		for propName, propSchema := range schema.Properties {
			tsType := g.openAPITypeToTS(propSchema)
			b.WriteString(fmt.Sprintf("  %s: %s;\n", propName, tsType))
		}
	}

	b.WriteString("}\n\n")

	return b.String()
}

// openAPITypeToTS converts OpenAPI type to TypeScript type
func (g *Generator) openAPITypeToTS(schema *exporter.OpenAPISchema) string {
	if schema == nil {
		return "any"
	}

	if schema.Ref != "" {
		// Extract type name from $ref
		parts := strings.Split(schema.Ref, "/")
		return parts[len(parts)-1]
	}

	switch schema.Type {
	case "string":
		if schema.Format == "date-time" || schema.Format == "date" {
			return "string" // Could use Date, but string is more compatible
		}
		return "string"
	case "integer", "number":
		return "number"
	case "boolean":
		return "boolean"
	case "array":
		if schema.Items != nil {
			itemType := g.openAPITypeToTS(schema.Items)
			return fmt.Sprintf("%s[]", itemType)
		}
		return "any[]"
	case "object":
		if schema.Properties != nil {
			// Inline object
			var props []string
			for name, prop := range schema.Properties {
				propType := g.openAPITypeToTS(prop)
				props = append(props, fmt.Sprintf("%s: %s", name, propType))
			}
			return fmt.Sprintf("{ %s }", strings.Join(props, "; "))
		}
		return "Record<string, any>"
	default:
		return "any"
	}
}

// generateClient generates the API client
func (g *Generator) generateClient() error {
	var b strings.Builder

	b.WriteString(`// Auto-generated API client from OpenAPI specification

import type * as Types from './types';

export interface ClientConfig {
  baseURL?: string;
  headers?: Record<string, string>;
  timeout?: number;
}

export class APIClient {
  private baseURL: string;
  private headers: Record<string, string>;
  private timeout: number;

  constructor(config: ClientConfig = {}) {
    this.baseURL = config.baseURL || '` + g.getDefaultBaseURL() + `';
    this.headers = config.headers || {};
    this.timeout = config.timeout || 30000;
  }

  private async request<T>(
    path: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = this.baseURL + path;
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), this.timeout);

    try {
      const response = await fetch(url, {
        ...options,
        headers: {
          'Content-Type': 'application/json',
          ...this.headers,
          ...options.headers,
        },
        signal: controller.signal,
      });

      clearTimeout(timeoutId);

      if (!response.ok) {
        const error: Types.APIError = {
          message: response.statusText,
          status: response.status,
        };

        try {
          const fault = await response.json();
          error.fault = fault;
        } catch {
          // No JSON body
        }

        throw error;
      }

      return await response.json();
    } catch (err) {
      clearTimeout(timeoutId);

      if (err instanceof Error && err.name === 'AbortError') {
        throw new Error('Request timeout');
      }

      throw err;
    }
  }

`)

	// Generate methods for each operation
	for path, pathItem := range g.spec.Paths {
		if pathItem.Post != nil {
			op := pathItem.Post
			methodName := toCamelCase(op.OperationID)
			requestType := toPascalCase(op.OperationID) + "Request"
			responseType := toPascalCase(op.OperationID) + "Response"

			b.WriteString(fmt.Sprintf("  /**\n   * %s\n", op.Summary))
			if op.Description != "" {
				b.WriteString(fmt.Sprintf("   * %s\n", op.Description))
			}
			b.WriteString("   */\n")
			b.WriteString(fmt.Sprintf("  async %s(request: Types.%s): Promise<Types.%s> {\n",
				methodName, requestType, responseType))
			b.WriteString(fmt.Sprintf("    return this.request<Types.%s>('%s', {\n", responseType, path))
			b.WriteString("      method: 'POST',\n")
			b.WriteString("      body: JSON.stringify(request),\n")
			b.WriteString("    });\n")
			b.WriteString("  }\n\n")
		}
	}

	b.WriteString("}\n\n")
	b.WriteString("// Export a default client instance\n")
	b.WriteString("export const apiClient = new APIClient();\n")

	return os.WriteFile(filepath.Join(g.outputDir, "client.ts"), []byte(b.String()), 0644)
}

// generateIndex generates the index file
func (g *Generator) generateIndex() error {
	content := `// Auto-generated API client exports

export * from './types';
export * from './client';
`

	return os.WriteFile(filepath.Join(g.outputDir, "index.ts"), []byte(content), 0644)
}

// generatePackageJSON generates package.json
func (g *Generator) generatePackageJSON() error {
	content := fmt.Sprintf(`{
  "name": "%s-client",
  "version": "1.0.0",
  "description": "TypeScript client for %s API",
  "main": "index.ts",
  "types": "index.ts",
  "scripts": {
    "build": "tsc",
    "type-check": "tsc --noEmit"
  },
  "keywords": ["api", "client", "typescript", "soap", "wsdl"],
  "author": "wsdl2api",
  "license": "MIT",
  "devDependencies": {
    "typescript": "^5.0.0"
  }
}
`, strings.ToLower(strings.ReplaceAll(g.spec.Info.Title, " ", "-")), g.spec.Info.Title)

	return os.WriteFile(filepath.Join(g.outputDir, "package.json"), []byte(content), 0644)
}

// getDefaultBaseURL gets the default base URL from servers
func (g *Generator) getDefaultBaseURL() string {
	if len(g.spec.Servers) > 0 {
		return g.spec.Servers[0].URL
	}
	return "http://localhost:8080"
}

// Helper functions
func toPascalCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	// Split by common separators
	words := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == '.' || r == ' '
	})

	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}

	return strings.Join(words, "")
}

func toCamelCase(s string) string {
	pascal := toPascalCase(s)
	if len(pascal) > 0 {
		return strings.ToLower(string(pascal[0])) + pascal[1:]
	}
	return pascal
}

// generateTSConfig generates tsconfig.json
func (g *Generator) generateTSConfig() error {
	content := `{
  "compilerOptions": {
    "target": "ES2020",
    "module": "ESNext",
    "moduleResolution": "bundler",
    "lib": ["ES2020", "DOM"],
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true,
    "declaration": true,
    "declarationMap": true,
    "sourceMap": true,
    "outDir": "./dist",
    "rootDir": "."
  },
  "include": ["*.ts"],
  "exclude": ["node_modules", "dist"]
}
`
	return os.WriteFile(filepath.Join(g.outputDir, "tsconfig.json"), []byte(content), 0644)
}

// generateReadme generates README.md
func (g *Generator) generateReadme() error {
	content := fmt.Sprintf(`# %s TypeScript Client

Auto-generated TypeScript client for %s API.

## Installation

` + "```bash\n" + `npm install
` + "```\n" + `

## Usage

` + "```typescript\n" + `import { APIClient } from './client';

// Create client instance
const client = new APIClient({
  baseURL: '%s',
  headers: {
    // Add custom headers if needed
  },
  timeout: 30000 // 30 seconds
});

// Make API calls
try {
  const response = await client.someOperation(request);
  console.log(response);
} catch (error) {
  console.error('API Error:', error);
}
` + "```\n" + `

## Type Safety

This client is fully typed with TypeScript. All request and response types are available:

` + "```typescript\n" + `import type { SomeOperationRequest, SomeOperationResponse } from './types';

const request: SomeOperationRequest = {
  // Your request data (autocomplete available!)
};

const response: SomeOperationResponse = await client.someOperation(request);
` + "```\n" + `

## Error Handling

The client throws typed errors:

` + "```typescript\n" + `import type { APIError } from './types';

try {
  await client.someOperation(request);
} catch (error) {
  const apiError = error as APIError;
  console.error('Status:', apiError.status);
  console.error('Message:', apiError.message);
  if (apiError.fault) {
    console.error('SOAP Fault:', apiError.fault);
  }
}
` + "```\n" + `

## Build

` + "```bash\n" + `npm run build
` + "```\n" + `

## Type Check

` + "```bash\n" + `npm run type-check
` + "```\n" + `

---

Generated by [wsdl2api](https://github.com/thdev01/wsdl2api)
`, g.spec.Info.Title, g.spec.Info.Title, g.getDefaultBaseURL())

	return os.WriteFile(filepath.Join(g.outputDir, "README.md"), []byte(content), 0644)
}
