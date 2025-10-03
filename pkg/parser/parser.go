package parser

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/thdev01/wsdl2api/internal/models"
)

// Parser handles WSDL parsing
type Parser struct{}

// NewParser creates a new WSDL parser
func NewParser() *Parser {
	return &Parser{}
}

// Parse parses a WSDL from file or URL
func (p *Parser) Parse(path string) (*models.Definitions, error) {
	var reader io.ReadCloser
	var err error

	// Check if path is URL or file
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		// Fetch from URL
		resp, err := http.Get(path)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch WSDL from URL: %w", err)
		}
		reader = resp.Body
	} else {
		// Read from file
		file, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("failed to open WSDL file: %w", err)
		}
		reader = file
	}
	defer reader.Close()

	// Parse XML
	var rawWSDL rawDefinitions
	decoder := xml.NewDecoder(reader)
	if err = decoder.Decode(&rawWSDL); err != nil {
		return nil, fmt.Errorf("failed to decode WSDL XML: %w", err)
	}

	// Convert to internal model
	return p.convertToModel(&rawWSDL), nil
}

// convertToModel converts raw XML structures to internal models
func (p *Parser) convertToModel(raw *rawDefinitions) *models.Definitions {
	def := &models.Definitions{
		Name:            raw.Name,
		TargetNamespace: raw.TargetNamespace,
		Services:        make([]models.Service, 0),
		Bindings:        make([]models.Binding, 0),
		PortTypes:       make([]models.PortType, 0),
		Messages:        make([]models.Message, 0),
	}

	// Convert services
	for _, svc := range raw.Service {
		service := models.Service{
			Name:  svc.Name,
			Ports: make([]models.Port, 0),
		}
		for _, port := range svc.Port {
			service.Ports = append(service.Ports, models.Port{
				Name:    port.Name,
				Binding: port.Binding,
				Address: port.Address.Location,
			})
		}
		def.Services = append(def.Services, service)
	}

	// Convert bindings
	for _, bind := range raw.Binding {
		binding := models.Binding{
			Name:       bind.Name,
			Type:       bind.Type,
			Operations: make([]models.BindingOperation, 0),
		}
		for _, op := range bind.Operation {
			operation := models.BindingOperation{
				Name:       op.Name,
				SoapAction: op.SoapOperation.SoapAction,
			}
			binding.Operations = append(binding.Operations, operation)
		}
		def.Bindings = append(def.Bindings, binding)
	}

	// Convert port types
	for _, pt := range raw.PortType {
		portType := models.PortType{
			Name:       pt.Name,
			Operations: make([]models.Operation, 0),
		}
		for _, op := range pt.Operation {
			operation := models.Operation{
				Name:          op.Name,
				Documentation: op.Documentation,
				Input: models.Message{
					Name: op.Input.Message,
				},
				Output: models.Message{
					Name: op.Output.Message,
				},
			}
			portType.Operations = append(portType.Operations, operation)
		}
		def.PortTypes = append(def.PortTypes, portType)
	}

	// Convert messages
	for _, msg := range raw.Message {
		message := models.Message{
			Name:  msg.Name,
			Parts: make([]models.Part, 0),
		}
		for _, part := range msg.Part {
			message.Parts = append(message.Parts, models.Part{
				Name:    part.Name,
				Element: part.Element,
				Type:    part.Type,
			})
		}
		def.Messages = append(def.Messages, message)
	}

	return def
}

// Raw XML structures for unmarshaling
type rawDefinitions struct {
	XMLName         xml.Name      `xml:"definitions"`
	Name            string        `xml:"name,attr"`
	TargetNamespace string        `xml:"targetNamespace,attr"`
	Service         []rawService  `xml:"service"`
	Binding         []rawBinding  `xml:"binding"`
	PortType        []rawPortType `xml:"portType"`
	Message         []rawMessage  `xml:"message"`
	Types           rawTypes      `xml:"types"`
}

type rawService struct {
	Name string    `xml:"name,attr"`
	Port []rawPort `xml:"port"`
}

type rawPort struct {
	Name    string     `xml:"name,attr"`
	Binding string     `xml:"binding,attr"`
	Address rawAddress `xml:"address"`
}

type rawAddress struct {
	Location string `xml:"location,attr"`
}

type rawBinding struct {
	Name      string             `xml:"name,attr"`
	Type      string             `xml:"type,attr"`
	Operation []rawBindOperation `xml:"operation"`
}

type rawBindOperation struct {
	Name          string         `xml:"name,attr"`
	SoapOperation rawSoapOperation `xml:"operation"`
	Input         rawBindMessage `xml:"input"`
	Output        rawBindMessage `xml:"output"`
}

type rawSoapOperation struct {
	SoapAction string `xml:"soapAction,attr"`
}

type rawBindMessage struct {
	Body rawBody `xml:"body"`
}

type rawBody struct {
	Use string `xml:"use,attr"`
}

type rawPortType struct {
	Name      string         `xml:"name,attr"`
	Operation []rawOperation `xml:"operation"`
}

type rawOperation struct {
	Name          string             `xml:"name,attr"`
	Documentation string             `xml:"documentation"`
	Input         rawOperationMessage `xml:"input"`
	Output        rawOperationMessage `xml:"output"`
}

type rawOperationMessage struct {
	Message string `xml:"message,attr"`
}

type rawMessage struct {
	Name string    `xml:"name,attr"`
	Part []rawPart `xml:"part"`
}

type rawPart struct {
	Name    string `xml:"name,attr"`
	Element string `xml:"element,attr"`
	Type    string `xml:"type,attr"`
}

type rawTypes struct {
	Schema []rawSchema `xml:"schema"`
}

type rawSchema struct {
	TargetNamespace string          `xml:"targetNamespace,attr"`
	Element         []rawXSDElement `xml:"element"`
}

type rawXSDElement struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}
