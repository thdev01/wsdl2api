# WSDL2API

**Convert legacy SOAP/WSDL services into modern REST APIs**

Transform WSDL definitions into clean, modular Go code with automatically generated REST API endpoints.

---

## Features

- 🔄 **WSDL Parsing**: Parse any WSDL file (local or remote)
- 🏗️ **Code Generation**: Generate complete Go client structures
- 🌐 **REST API**: Automatically create RESTful endpoints for SOAP operations
- 🔐 **WS-Security**: Full authentication support (UsernameToken, Digest)
- 📊 **Complex Types**: Handle nested structs, arrays, and optional fields
- 🔧 **SOAP 1.1 & 1.2**: Support for both SOAP protocol versions
- 🎭 **Mock Server**: Generate mock SOAP servers for testing
- 📄 **OpenAPI Export**: Convert WSDL to OpenAPI 3.0 specifications
- 💙 **TypeScript Client**: Generate type-safe TypeScript/JavaScript clients
- 📦 **Modular**: Clean, organized code structure
- 🚀 **Easy to Use**: Simple CLI interface
- 🧪 **Tested**: Pre-tested with real-world WSDLs (Correios, etc.)

---

## Installation

```bash
go install github.com/thdev01/wsdl2api/cmd/wsdl2api@latest
```

Or build from source:

```bash
git clone https://github.com/thdev01/wsdl2api.git
cd wsdl2api
go build -o wsdl2api ./cmd/wsdl2api
```

---

## Quick Start

### Generate Client from WSDL

```bash
# Basic generation
wsdl2api generate --wsdl https://example.com/service?wsdl --output ./generated

# With SOAP 1.2 support
wsdl2api generate --wsdl ./service.wsdl --soap-version 1.2 --output ./generated

# With mock server for testing
wsdl2api generate --wsdl ./service.wsdl --mock --output ./generated
```

### Export to OpenAPI & TypeScript

```bash
# Export to JSON
wsdl2api export --wsdl ./service.wsdl --format json --output ./docs

# Export to YAML
wsdl2api export --wsdl ./service.wsdl --format yaml --output ./docs

# Generate TypeScript client alongside OpenAPI
wsdl2api export --wsdl ./service.wsdl --output ./api --typescript

# Custom TypeScript output directory
wsdl2api export --wsdl ./service.wsdl --output ./api --typescript --ts-output ./client
```

### Start REST API Server

```bash
# Generate and start server
wsdl2api serve --wsdl https://example.com/service?wsdl --port 8080
```

### Example: Correios CEP Service

```bash
# Brazilian postal code lookup service
wsdl2api serve --wsdl https://apps.correios.com.br/SigepMasterJPA/AtendeClienteService/AtendeCliente?wsdl --port 8080

# Test the endpoint
curl http://localhost:8080/api/consultaCEP -X POST -H "Content-Type: application/json" -d '{"cep": "01310-100"}'
```

---

## Usage

### Generate Go Client Code

```bash
wsdl2api generate \
  --wsdl <wsdl-url-or-path> \
  --output <output-directory> \
  --package <package-name>
```

#### Generated Files:
- `client.go` - SOAP client with WS-Security and SOAP 1.1/1.2 support
- `types.go` - Request/response types with complex type handling
- `operators.go` - Easy-to-use functions for each operation
- `example.go` - Usage documentation
- `mock_server.go` - Mock server for testing (with --mock flag)

#### Use Generated Code:

```go
package main

import (
    "fmt"
    "log"
    "yourproject/generated/client"
)

func main() {
    // Create client
    c := client.NewClient("")

    // Optional: Set WS-Security authentication
    c.SetBasicAuth("username", "password")
    // Or use digest authentication
    // c.SetDigestAuth("username", "password")

    // Optional: Use SOAP 1.2
    // c.SetSOAPVersion("1.2")

    // Call operation with seamless API
    result, err := c.SomeOperation(param1, param2)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Result: %+v\n", result)
}
```

#### Use Generated TypeScript Client:

```typescript
import { APIClient } from './typescript/client';
import type { AddRequest, AddResponse } from './typescript/types';

// Create client instance
const client = new APIClient({
  baseURL: 'http://your-api-url.com',
  timeout: 30000,
  headers: {
    'Authorization': 'Bearer token'  // Optional
  }
});

// Make type-safe API calls
const request: AddRequest = {
  parameters: 'value'
};

try {
  const response: AddResponse = await client.add(request);
  console.log('Response:', response);
} catch (error) {
  const apiError = error as APIError;
  console.error('Error:', apiError.message);
  if (apiError.fault) {
    console.error('SOAP Fault:', apiError.fault);
  }
}
```

### Serve REST API

```bash
wsdl2api serve \
  --wsdl <wsdl-url-or-path> \
  --port <port-number> \
  --host <host-address>
```

### Options

#### Generate Command
```
Flags:
  -w, --wsdl string        WSDL file path or URL (required)
  -o, --output string      Output directory (default "./generated")
  -p, --package string     Go package name (default "client")
  --mock                   Generate mock server for testing
  --soap-version string    SOAP version: "1.1" or "1.2" (default "1.1")
  -h, --help              Help for command
```

#### Export Command
```
Flags:
  -w, --wsdl string        WSDL file path or URL (required)
  -o, --output string      Output directory (empty for stdout)
  -f, --format string      Export format: "json" or "yaml" (default "json")
  --typescript             Generate TypeScript client
  --ts-output string       TypeScript output directory (default: <output>/typescript)
  -h, --help              Help for command
```

#### Serve Command
```
Flags:
  -w, --wsdl string    WSDL file path or URL (required)
  --port int          Server port (default 8080)
  --host string       Server host (default "localhost")
  -h, --help          Help for command
```

📚 **[Complete Usage Guide](docs/USAGE.md)** - Advanced examples, best practices, troubleshooting

---

## Project Structure

```
wsdl2api/
├── cmd/
│   └── wsdl2api/          # CLI application
├── pkg/
│   ├── parser/            # WSDL parsing logic
│   ├── generator/         # Code generation (client, types, operators, mock)
│   ├── security/          # WS-Security implementation
│   ├── exporter/          # OpenAPI/Swagger export
│   ├── typescript/        # TypeScript client generator
│   ├── client/            # SOAP client wrapper
│   └── server/            # REST API server
├── internal/
│   ├── models/            # Data models
│   └── utils/             # Utilities
├── examples/              # Example WSDLs and usage
├── docs/                  # Documentation
└── tests/                 # Test suite
```

---

## How It Works

1. **Parse WSDL**: Extract services, operations, and data types
2. **Generate Structs**: Create Go structs for all WSDL types
3. **Create Client**: Generate SOAP client code
4. **Build API**: Create REST endpoints for each SOAP operation
5. **Serve**: Run HTTP server with generated routes

---

## Examples

See [examples/](examples/) directory for:
- Correios (Brazilian Postal Service)
- Public SOAP services
- Custom WSDL examples

---

## Development

```bash
# Install dependencies
go mod download

# Run tests
go test ./...

# Build
go build -o wsdl2api ./cmd/wsdl2api

# Run
./wsdl2api serve --wsdl examples/correios.wsdl
```

---

## Roadmap

- [x] WSDL parsing
- [x] Code generation
- [x] REST API generation
- [x] CLI interface
- [x] Support for complex types (nested structs, arrays)
- [x] WS-Security authentication support
- [x] SOAP 1.2 support
- [x] Mock server generation for testing
- [x] OpenAPI/Swagger 3.0 export
- [x] TypeScript/JavaScript client generation
- [ ] Docker container
- [ ] Web UI
- [ ] Advanced type validation
- [ ] Custom headers support
- [ ] Python client generation
- [ ] GraphQL API generation

---

## Contributing

Contributions welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details.

---

## License

MIT License - see [LICENSE](LICENSE) for details.

---

## Author

**thdev01** (thdev01@gmail.com)
