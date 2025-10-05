package server

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/thdev01/wsdl2api/internal/models"
)

// Server represents the REST API server
type Server struct {
	definitions  *models.Definitions
	host         string
	port         int
	router       *gin.Engine
	soapEndpoint string
	soapVersion  string
}

// NewServer creates a new REST API server
func NewServer(def *models.Definitions, host string, port int) *Server {
	// Extract SOAP endpoint from definitions
	soapEndpoint := ""
	if len(def.Services) > 0 && len(def.Services[0].Ports) > 0 {
		soapEndpoint = def.Services[0].Ports[0].Address
	}

	return &Server{
		definitions:  def,
		host:         host,
		port:         port,
		router:       gin.Default(),
		soapEndpoint: soapEndpoint,
		soapVersion:  "1.1", // Default to SOAP 1.1
	}
}

// SetSOAPEndpoint sets a custom SOAP endpoint
func (s *Server) SetSOAPEndpoint(endpoint string) {
	s.soapEndpoint = endpoint
}

// SetSOAPVersion sets the SOAP version (1.1 or 1.2)
func (s *Server) SetSOAPVersion(version string) {
	s.soapVersion = version
}

// Start starts the REST API server
func (s *Server) Start() error {
	// Setup routes
	s.setupRoutes()

	// Start server
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	return s.router.Run(addr)
}

// setupRoutes configures all API routes
func (s *Server) setupRoutes() {
	// Health check
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"service": s.definitions.Name,
		})
	})

	// Service info
	s.router.GET("/info", s.handleServiceInfo)

	// API routes group
	api := s.router.Group("/api")

	// Generate routes for each operation in each port type
	for _, portType := range s.definitions.PortTypes {
		for _, op := range portType.Operations {
			// Create REST endpoint for SOAP operation
			path := fmt.Sprintf("/%s", op.Name)
			api.POST(path, s.createOperationHandler(op))
			api.GET(path+"/info", s.createOperationInfoHandler(op))
		}
	}
}

// handleServiceInfo returns service information
func (s *Server) handleServiceInfo(c *gin.Context) {
	services := make([]gin.H, 0)
	for _, svc := range s.definitions.Services {
		ports := make([]gin.H, 0)
		for _, port := range svc.Ports {
			ports = append(ports, gin.H{
				"name":    port.Name,
				"binding": port.Binding,
				"address": port.Address,
			})
		}
		services = append(services, gin.H{
			"name":  svc.Name,
			"ports": ports,
		})
	}

	operations := make([]gin.H, 0)
	for _, portType := range s.definitions.PortTypes {
		for _, op := range portType.Operations {
			operations = append(operations, gin.H{
				"name":          op.Name,
				"documentation": op.Documentation,
				"endpoint":      fmt.Sprintf("/api/%s", op.Name),
				"method":        "POST",
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"name":            s.definitions.Name,
		"targetNamespace": s.definitions.TargetNamespace,
		"services":        services,
		"operations":      operations,
		"totalOperations": len(operations),
	})
}

// createOperationHandler creates a handler for a SOAP operation
func (s *Server) createOperationHandler(op models.Operation) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse request body
		var requestBody map[string]interface{}
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request body",
				"details": err.Error(),
			})
			return
		}

		// Find SOAP action for this operation
		soapAction := ""
		for _, binding := range s.definitions.Bindings {
			for _, bindOp := range binding.Operations {
				if bindOp.Name == op.Name {
					soapAction = bindOp.SoapAction
					break
				}
			}
		}

		// Make actual SOAP call
		response, err := s.callSOAP(op.Name, soapAction, requestBody)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":     "SOAP call failed",
				"operation": op.Name,
				"details":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"operation": op.Name,
			"status":    "success",
			"request":   requestBody,
			"response":  response,
		})
	}
}

// createOperationInfoHandler creates an info handler for an operation
func (s *Server) createOperationInfoHandler(op models.Operation) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Find SOAP action
		soapAction := ""
		for _, binding := range s.definitions.Bindings {
			for _, bindOp := range binding.Operations {
				if bindOp.Name == op.Name {
					soapAction = bindOp.SoapAction
					break
				}
			}
		}

		// Find message details
		inputParts := make([]gin.H, 0)
		outputParts := make([]gin.H, 0)

		for _, msg := range s.definitions.Messages {
			if msg.Name == op.Input.Name {
				for _, part := range msg.Parts {
					inputParts = append(inputParts, gin.H{
						"name":    part.Name,
						"type":    part.Type,
						"element": part.Element,
					})
				}
			}
			if msg.Name == op.Output.Name {
				for _, part := range msg.Parts {
					outputParts = append(outputParts, gin.H{
						"name":    part.Name,
						"type":    part.Type,
						"element": part.Element,
					})
				}
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"operation":     op.Name,
			"documentation": op.Documentation,
			"soapAction":    soapAction,
			"endpoint":      fmt.Sprintf("/api/%s", op.Name),
			"method":        "POST",
			"input": gin.H{
				"message": op.Input.Name,
				"parts":   inputParts,
			},
			"output": gin.H{
				"message": op.Output.Name,
				"parts":   outputParts,
			},
			"example": gin.H{
				"curl": fmt.Sprintf(`curl -X POST http://%s:%d/api/%s \
  -H "Content-Type: application/json" \
  -d '{"param": "value"}'`, s.host, s.port, op.Name),
			},
		})
	}
}

// callSOAP makes an actual SOAP call to the backend service
func (s *Server) callSOAP(operation, soapAction string, requestParams map[string]interface{}) (map[string]interface{}, error) {
	if s.soapEndpoint == "" {
		return nil, fmt.Errorf("SOAP endpoint not configured")
	}

	// Build SOAP envelope (returns XML string)
	xmlData := s.buildSOAPEnvelope(operation, requestParams)

	// Create HTTP request
	req, err := http.NewRequest("POST", s.soapEndpoint, bytes.NewBuffer([]byte(xmlData)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers based on SOAP version
	if s.soapVersion == "1.2" {
		req.Header.Set("Content-Type", "application/soap+xml; charset=utf-8")
	} else {
		req.Header.Set("Content-Type", "text/xml; charset=utf-8")
		if soapAction != "" {
			req.Header.Set("SOAPAction", fmt.Sprintf(`"%s"`, soapAction))
		}
	}

	// Make the call
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("SOAP call failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse SOAP response
	result, err := s.parseSOAPResponse(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SOAP response: %w", err)
	}

	return result, nil
}

// buildSOAPEnvelope builds a SOAP envelope for the request
func (s *Server) buildSOAPEnvelope(operation string, params map[string]interface{}) string {
	// Build parameter XML elements
	var paramsXML strings.Builder
	for k, v := range params {
		paramsXML.WriteString(fmt.Sprintf("<%s>%v</%s>", k, v, k))
	}

	// Get target namespace from definitions
	targetNS := s.definitions.TargetNamespace
	if targetNS == "" {
		targetNS = "http://tempuri.org/"
	}

	if s.soapVersion == "1.2" {
		return fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<soap12:Envelope xmlns:soap12="http://www.w3.org/2003/05/soap-envelope" xmlns:tns="%s">
  <soap12:Body>
    <tns:%s>%s</tns:%s>
  </soap12:Body>
</soap12:Envelope>`, targetNS, operation, paramsXML.String(), operation)
	}

	// SOAP 1.1
	return fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:tns="%s">
  <soap:Body>
    <tns:%s>%s</tns:%s>
  </soap:Body>
</soap:Envelope>`, targetNS, operation, paramsXML.String(), operation)
}

// parseSOAPResponse parses a SOAP response and extracts the result
func (s *Server) parseSOAPResponse(xmlData []byte) (map[string]interface{}, error) {
	// Generic SOAP envelope structure
	var envelope struct {
		Body struct {
			Content string `xml:",innerxml"`
		} `xml:"Body"`
	}

	if err := xml.Unmarshal(xmlData, &envelope); err != nil {
		return nil, fmt.Errorf("failed to unmarshal SOAP envelope: %w", err)
	}

	// Try to parse the body content as JSON-friendly map
	result := make(map[string]interface{})

	// Simple XML to map conversion (can be enhanced)
	bodyContent := strings.TrimSpace(envelope.Body.Content)
	if bodyContent != "" {
		// For now, return the raw XML in the response
		result["xml"] = bodyContent

		// Try to extract values (basic implementation)
		result["raw"] = bodyContent
	}

	return result, nil
}
