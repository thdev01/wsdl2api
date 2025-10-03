package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thdev01/wsdl2api/internal/models"
)

// Server represents the REST API server
type Server struct {
	definitions *models.Definitions
	host        string
	port        int
	router      *gin.Engine
}

// NewServer creates a new REST API server
func NewServer(def *models.Definitions, host string, port int) *Server {
	return &Server{
		definitions: def,
		host:        host,
		port:        port,
		router:      gin.Default(),
	}
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
				"error": "Invalid request body",
				"details": err.Error(),
			})
			return
		}

		// TODO: Implement actual SOAP call
		// For now, return a mock response
		c.JSON(http.StatusOK, gin.H{
			"operation": op.Name,
			"status":    "success",
			"message":   "Operation executed successfully (mock response)",
			"request":   requestBody,
			"note":      "This is a mock response. SOAP client integration coming soon.",
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
