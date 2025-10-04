# WSDL2API Usage Guide

Complete guide for using wsdl2api to convert SOAP services into Go clients.

---

## Quick Start

### 1. Generate Client Code

```bash
# From WSDL URL
wsdl2api generate --wsdl http://example.com/service?wsdl --output ./client --package soapclient

# From local WSDL file
wsdl2api generate --wsdl service.wsdl --output ./client --package soapclient
```

### 2. Use Generated Code

```go
package main

import (
    "fmt"
    "log"

    "yourproject/client"
)

func main() {
    // Create client
    client := soapclient.NewClient("")

    // Call operations
    result, err := client.SomeOperation(param1, param2)
    if err != nil {
        log.Fatalf("Operation failed: %v", err)
    }

    fmt.Printf("Result: %+v\n", result)
}
```

---

## Generated Files

When you run `wsdl2api generate`, it creates:

```
output-dir/
├── client.go      # SOAP client with HTTP handling
├── types.go       # Request/response types
├── operators.go   # Easy-to-use operation functions
└── example.go     # Usage examples and documentation
```

### client.go

Contains the core SOAP client:

```go
type Client struct {
    URL        string
    HTTPClient *http.Client
    Headers    map[string]string
}

func NewClient(url string) *Client
func (c *Client) Call(soapAction string, request, response interface{}) error
func (c *Client) SetHeader(key, value string)
```

**Features:**
- Proper SOAP envelope handling
- XML marshaling/unmarshaling
- Custom HTTP headers support
- Error handling

### types.go

Request and response types for each operation:

```go
// Example for "Add" operation
type AddRequest struct {
    XMLName xml.Name `xml:"http://tempuri.org/ Add"`
    IntA    int      `xml:"intA"`
    IntB    int      `xml:"intB"`
}

type AddResponse struct {
    XMLName   xml.Name `xml:"AddResponse"`
    AddResult int      `xml:"AddResult"`
}
```

### operators.go

High-level functions for easy usage:

```go
// Add is an easy-to-use operator for the Add operation
func (c *Client) Add(intA int, intB int) (int, error) {
    request := &AddRequest{IntA: intA, IntB: intB}
    var response AddResponse

    err := c.Call("http://tempuri.org/Add", request, &response)
    if err != nil {
        return 0, fmt.Errorf("failed to execute Add: %w", err)
    }

    return response.AddResult, nil
}
```

---

## Usage Examples

### Basic Usage

```go
package main

import (
    "fmt"
    "log"

    "yourproject/generated/calculator"
)

func main() {
    // Create client (uses default URL from WSDL)
    client := calculator.NewClient("")

    // Call Add operation
    result, err := client.Add(5, 3)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("5 + 3 = %d\n", result)
}
```

### Custom URL

```go
// Override service URL
client := calculator.NewClient("http://my-soap-service.com/calculator.asmx")

result, err := client.Add(10, 20)
```

### Custom Headers

```go
client := calculator.NewClient("")

// Add custom headers
client.SetHeader("X-API-Key", "your-key")
client.SetHeader("Authorization", "Bearer token")

result, err := client.Add(5, 3)
```

### Error Handling

```go
result, err := client.Add(5, 3)
if err != nil {
    // Handle different error types
    switch {
    case strings.Contains(err.Error(), "failed to marshal"):
        log.Println("Invalid request structure")
    case strings.Contains(err.Error(), "failed to execute"):
        log.Println("Network error")
    case strings.Contains(err.Error(), "failed to unmarshal"):
        log.Println("Invalid response")
    default:
        log.Printf("Unknown error: %v", err)
    }
    return
}
```

---

## Advanced Usage

### Custom HTTP Client

```go
import (
    "net/http"
    "time"
)

client := calculator.NewClient("")

// Configure custom HTTP client
client.HTTPClient = &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:       10,
        IdleConnTimeout:    30 * time.Second,
        DisableCompression: false,
    },
}
```

### Logging Requests

```go
import (
    "bytes"
    "io"
    "log"
)

// Wrap HTTP client to log requests
type LoggingTransport struct {
    Transport http.RoundTripper
}

func (t *LoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    // Log request
    body, _ := io.ReadAll(req.Body)
    log.Printf("SOAP Request: %s", string(body))
    req.Body = io.NopCloser(bytes.NewBuffer(body))

    // Execute
    resp, err := t.Transport.RoundTrip(req)

    // Log response
    if resp != nil {
        respBody, _ := io.ReadAll(resp.Body)
        log.Printf("SOAP Response: %s", string(respBody))
        resp.Body = io.NopCloser(bytes.NewBuffer(respBody))
    }

    return resp, err
}

// Use it
client := calculator.NewClient("")
client.HTTPClient.Transport = &LoggingTransport{
    Transport: http.DefaultTransport,
}
```

### Authentication

```go
// Basic Auth
client := calculator.NewClient("")
client.SetHeader("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("username:password")))

// Bearer Token
client.SetHeader("Authorization", "Bearer your-jwt-token")

// API Key
client.SetHeader("X-API-Key", "your-api-key")
```

---

## REST API Server Mode

Instead of generating code, you can run a REST API proxy:

```bash
wsdl2api serve --wsdl service.wsdl --port 8080
```

This creates REST endpoints for each SOAP operation:

```bash
# Get service info
curl http://localhost:8080/info

# Call operation
curl -X POST http://localhost:8080/api/Add \
  -H "Content-Type: application/json" \
  -d '{"intA": 5, "intB": 3}'

# Get operation info
curl http://localhost:8080/api/Add/info
```

---

## Best Practices

### 1. Use Operators

✅ **Do:**
```go
result, err := client.Add(5, 3)
```

❌ **Don't:**
```go
req := &AddRequest{IntA: 5, IntB: 3}
var resp AddResponse
err := client.Call("...", req, &resp)
```

### 2. Handle Errors Properly

```go
result, err := client.SomeOperation(params)
if err != nil {
    // Log with context
    log.Printf("Operation failed: %v", err)
    // Return meaningful error to caller
    return fmt.Errorf("failed to process: %w", err)
}
```

### 3. Reuse Client Instances

```go
// Create once
var soapClient = calculator.NewClient("")

// Use many times
func processData(a, b int) error {
    result, err := soapClient.Add(a, b)
    // ...
}
```

### 4. Set Timeouts

```go
client := calculator.NewClient("")
client.HTTPClient.Timeout = 10 * time.Second
```

---

## Troubleshooting

### Connection Errors

```
failed to execute request: connection refused
```

**Solutions:**
- Check service URL is correct
- Verify service is accessible
- Check firewall/network settings

### Unmarshal Errors

```
failed to unmarshal response: XML syntax error
```

**Solutions:**
- Check SOAP envelope structure
- Verify namespace in XML tags
- Enable logging to inspect raw response

### SOAP Faults

SOAP services may return faults instead of errors. Check the response:

```go
result, err := client.SomeOperation(params)
if err != nil {
    if strings.Contains(err.Error(), "faultcode") {
        log.Println("SOAP Fault:", err)
    }
}
```

---

## Examples

See [examples/](../examples/) directory for:
- [use_calculator/](../examples/use_calculator/) - Calculator service example
- [WSDL_CATALOG.md](../examples/WSDL_CATALOG.md) - Public WSDLs for testing

---

## CLI Reference

```bash
# Generate code
wsdl2api generate [flags]

Flags:
  -w, --wsdl string      WSDL file path or URL (required)
  -o, --output string    Output directory (default "./generated")
  -p, --package string   Go package name (default "client")

# Serve REST API
wsdl2api serve [flags]

Flags:
  -w, --wsdl string     WSDL file path or URL (required)
  --port int           Server port (default 8080)
  --host string        Server host (default "localhost")
```

---

## Contributing

Found a bug or want to contribute? See [CONTRIBUTING.md](../CONTRIBUTING.md)
