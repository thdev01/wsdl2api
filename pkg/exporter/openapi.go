package exporter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/thdev01/wsdl2api/internal/models"
)

// OpenAPISpec represents an OpenAPI 3.0 specification
type OpenAPISpec struct {
	OpenAPI string                 `json:"openapi"`
	Info    OpenAPIInfo            `json:"info"`
	Servers []OpenAPIServer        `json:"servers,omitempty"`
	Paths   map[string]OpenAPIPath `json:"paths"`
	Components *OpenAPIComponents  `json:"components,omitempty"`
}

// OpenAPIInfo contains API metadata
type OpenAPIInfo struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Version     string `json:"version"`
}

// OpenAPIServer describes a server
type OpenAPIServer struct {
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}

// OpenAPIPath describes operations on a path
type OpenAPIPath struct {
	Post *OpenAPIOperation `json:"post,omitempty"`
	Get  *OpenAPIOperation `json:"get,omitempty"`
}

// OpenAPIOperation describes a single operation
type OpenAPIOperation struct {
	Summary     string                        `json:"summary,omitempty"`
	Description string                        `json:"description,omitempty"`
	OperationID string                        `json:"operationId,omitempty"`
	RequestBody *OpenAPIRequestBody           `json:"requestBody,omitempty"`
	Responses   map[string]OpenAPIResponse    `json:"responses"`
	Tags        []string                      `json:"tags,omitempty"`
}

// OpenAPIRequestBody describes a request body
type OpenAPIRequestBody struct {
	Description string                      `json:"description,omitempty"`
	Required    bool                        `json:"required,omitempty"`
	Content     map[string]OpenAPIMediaType `json:"content"`
}

// OpenAPIResponse describes a response
type OpenAPIResponse struct {
	Description string                      `json:"description"`
	Content     map[string]OpenAPIMediaType `json:"content,omitempty"`
}

// OpenAPIMediaType describes a media type
type OpenAPIMediaType struct {
	Schema *OpenAPISchema `json:"schema,omitempty"`
}

// OpenAPISchema describes a schema
type OpenAPISchema struct {
	Type       string                    `json:"type,omitempty"`
	Properties map[string]*OpenAPISchema `json:"properties,omitempty"`
	Items      *OpenAPISchema            `json:"items,omitempty"`
	Ref        string                    `json:"$ref,omitempty"`
	Format     string                    `json:"format,omitempty"`
}

// OpenAPIComponents contains reusable components
type OpenAPIComponents struct {
	Schemas map[string]*OpenAPISchema `json:"schemas,omitempty"`
}

// ConvertWSDLToOpenAPI converts WSDL definitions to OpenAPI spec
func ConvertWSDLToOpenAPI(def *models.Definitions) (*OpenAPISpec, error) {
	spec := &OpenAPISpec{
		OpenAPI: "3.0.0",
		Info: OpenAPIInfo{
			Title:       def.Name,
			Description: fmt.Sprintf("API converted from WSDL: %s", def.TargetNamespace),
			Version:     "1.0.0",
		},
		Paths: make(map[string]OpenAPIPath),
		Components: &OpenAPIComponents{
			Schemas: make(map[string]*OpenAPISchema),
		},
	}

	// Add servers
	for _, svc := range def.Services {
		for _, port := range svc.Ports {
			if port.Address != "" {
				spec.Servers = append(spec.Servers, OpenAPIServer{
					URL:         port.Address,
					Description: fmt.Sprintf("%s - %s", svc.Name, port.Name),
				})
			}
		}
	}

	// Convert operations
	for _, portType := range def.PortTypes {
		for _, op := range portType.Operations {
			path := fmt.Sprintf("/api/%s", op.Name)

			// Find input/output messages
			inputMsg := findMessage(def, op.Input.Name)
			outputMsg := findMessage(def, op.Output.Name)

			operation := &OpenAPIOperation{
				Summary:     op.Name,
				Description: op.Documentation,
				OperationID: op.Name,
				Responses:   make(map[string]OpenAPIResponse),
			}

			// Add request body
			if inputMsg != nil {
				operation.RequestBody = &OpenAPIRequestBody{
					Description: fmt.Sprintf("Request for %s operation", op.Name),
					Required:    true,
					Content: map[string]OpenAPIMediaType{
						"application/json": {
							Schema: convertMessageToSchema(inputMsg),
						},
					},
				}
			}

			// Add response
			if outputMsg != nil {
				operation.Responses["200"] = OpenAPIResponse{
					Description: fmt.Sprintf("Successful response for %s", op.Name),
					Content: map[string]OpenAPIMediaType{
						"application/json": {
							Schema: convertMessageToSchema(outputMsg),
						},
					},
				}
			}

			// Add error response
			operation.Responses["500"] = OpenAPIResponse{
				Description: "SOAP Fault",
				Content: map[string]OpenAPIMediaType{
					"application/json": {
						Schema: &OpenAPISchema{
							Type: "object",
							Properties: map[string]*OpenAPISchema{
								"faultcode":   {Type: "string"},
								"faultstring": {Type: "string"},
								"detail":      {Type: "string"},
							},
						},
					},
				},
			}

			spec.Paths[path] = OpenAPIPath{
				Post: operation,
			}
		}
	}

	return spec, nil
}

// convertMessageToSchema converts a WSDL message to OpenAPI schema
func convertMessageToSchema(msg *models.Message) *OpenAPISchema {
	if len(msg.Parts) == 0 {
		return &OpenAPISchema{Type: "object"}
	}

	schema := &OpenAPISchema{
		Type:       "object",
		Properties: make(map[string]*OpenAPISchema),
	}

	for _, part := range msg.Parts {
		schema.Properties[part.Name] = xsdTypeToOpenAPISchema(part.Type)
	}

	return schema
}

// xsdTypeToOpenAPISchema converts XSD type to OpenAPI schema
func xsdTypeToOpenAPISchema(xsdType string) *OpenAPISchema {
	// Remove namespace prefix
	if idx := strings.LastIndex(xsdType, ":"); idx != -1 {
		xsdType = xsdType[idx+1:]
	}

	typeMap := map[string]OpenAPISchema{
		"string":   {Type: "string"},
		"int":      {Type: "integer", Format: "int32"},
		"integer":  {Type: "integer", Format: "int32"},
		"long":     {Type: "integer", Format: "int64"},
		"short":    {Type: "integer", Format: "int32"},
		"boolean":  {Type: "boolean"},
		"float":    {Type: "number", Format: "float"},
		"double":   {Type: "number", Format: "double"},
		"decimal":  {Type: "number"},
		"dateTime": {Type: "string", Format: "date-time"},
		"date":     {Type: "string", Format: "date"},
		"time":     {Type: "string", Format: "time"},
	}

	if schema, ok := typeMap[xsdType]; ok {
		return &schema
	}

	// Default to string for unknown types
	return &OpenAPISchema{Type: "string"}
}

// findMessage finds a message by name
func findMessage(def *models.Definitions, name string) *models.Message {
	// Remove namespace prefix
	if idx := strings.LastIndex(name, ":"); idx != -1 {
		name = name[idx+1:]
	}

	for _, msg := range def.Messages {
		msgName := msg.Name
		if idx := strings.LastIndex(msgName, ":"); idx != -1 {
			msgName = msgName[idx+1:]
		}
		if msgName == name {
			return &msg
		}
	}
	return nil
}

// ExportToJSON exports OpenAPI spec as JSON
func (spec *OpenAPISpec) ExportToJSON() (string, error) {
	data, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ExportToYAML exports OpenAPI spec as YAML (simplified)
func (spec *OpenAPISpec) ExportToYAML() (string, error) {
	// For now, export as JSON (full YAML support would require yaml package)
	json, err := spec.ExportToJSON()
	if err != nil {
		return "", err
	}
	return "# OpenAPI YAML export (use a YAML converter for proper formatting)\n" + json, nil
}
