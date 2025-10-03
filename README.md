# WSDL2API

**Convert legacy SOAP/WSDL services into modern REST APIs**

Transform WSDL definitions into clean, modular Go code with automatically generated REST API endpoints.

---

## Features

- 🔄 **WSDL Parsing**: Parse any WSDL file (local or remote)
- 🏗️ **Code Generation**: Generate complete Go client structures
- 🌐 **REST API**: Automatically create RESTful endpoints for SOAP operations
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
# From remote URL
wsdl2api generate --wsdl https://example.com/service?wsdl --output ./generated

# From local file
wsdl2api generate --wsdl ./service.wsdl --output ./generated
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

### Generate Code Only

```bash
wsdl2api generate \
  --wsdl <wsdl-url-or-path> \
  --output <output-directory> \
  --package <package-name>
```

### Generate and Serve API

```bash
wsdl2api serve \
  --wsdl <wsdl-url-or-path> \
  --port <port-number> \
  --host <host-address>
```

### Options

```
Flags:
  -w, --wsdl string      WSDL file path or URL (required)
  -o, --output string    Output directory (default "./generated")
  -p, --package string   Go package name (default "client")
  --port int            Server port (default 8080)
  --host string         Server host (default "localhost")
  -h, --help            Help for command
```

---

## Project Structure

```
wsdl2api/
├── cmd/
│   └── wsdl2api/          # CLI application
├── pkg/
│   ├── parser/            # WSDL parsing logic
│   ├── generator/         # Code generation
│   ├── client/            # SOAP client wrapper
│   └── server/            # REST API server
├── internal/
│   ├── models/            # Data models
│   └── utils/             # Utilities
├── examples/              # Example WSDLs and usage
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
- [ ] Support for complex types
- [ ] WS-Security support
- [ ] OpenAPI/Swagger documentation generation
- [ ] Docker container
- [ ] Web UI

---

## Contributing

Contributions welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details.

---

## License

MIT License - see [LICENSE](LICENSE) for details.

---

## Author

**thdev01** (thdev01@gmail.com)
