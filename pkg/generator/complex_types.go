package generator

import (
	"fmt"
	"strings"

	"github.com/thdev01/wsdl2api/internal/models"
)

// ComplexTypeGenerator handles complex type generation
type ComplexTypeGenerator struct {
	targetNamespace string
	generatedTypes  map[string]bool
}

// NewComplexTypeGenerator creates a new complex type generator
func NewComplexTypeGenerator(targetNS string) *ComplexTypeGenerator {
	return &ComplexTypeGenerator{
		targetNamespace: targetNS,
		generatedTypes:  make(map[string]bool),
	}
}

// GenerateComplexType generates Go code for a complex type
func (ctg *ComplexTypeGenerator) GenerateComplexType(t models.Type) string {
	if ctg.generatedTypes[t.Name] {
		return ""
	}

	var b strings.Builder
	typeName := toPascalCase(t.Name)

	b.WriteString(fmt.Sprintf("// %s represents a complex type from WSDL\n", typeName))
	b.WriteString(fmt.Sprintf("type %s struct {\n", typeName))
	b.WriteString(fmt.Sprintf("\tXMLName xml.Name `xml:\"%s %s\"`\n", ctg.targetNamespace, t.Name))

	// Generate fields for elements
	for _, elem := range t.Elements {
		fieldName := toPascalCase(elem.Name)
		fieldType := ctg.getFieldType(elem)
		xmlTag := ctg.buildXMLTag(elem)

		b.WriteString(fmt.Sprintf("\t%s %s `xml:\"%s\"`\n", fieldName, fieldType, xmlTag))
	}

	// Generate fields for attributes
	for _, attr := range t.Attributes {
		fieldName := toPascalCase(attr.Name)
		fieldType := mapXSDTypeToGo(attr.Type)

		b.WriteString(fmt.Sprintf("\t%s %s `xml:\"%s,attr\"`\n", fieldName, fieldType, attr.Name))
	}

	b.WriteString("}\n\n")

	ctg.generatedTypes[t.Name] = true
	return b.String()
}

// getFieldType determines the Go type for an element
func (ctg *ComplexTypeGenerator) getFieldType(elem models.Element) string {
	baseType := mapXSDTypeToGo(elem.Type)

	// Handle arrays (maxOccurs > 1 or "unbounded")
	if elem.MaxOccurs == "unbounded" || (elem.MaxOccurs != "" && elem.MaxOccurs != "1") {
		baseType = "[]" + baseType
	}

	// Handle optional fields (minOccurs = 0)
	if elem.MinOccurs == "0" && !strings.HasPrefix(baseType, "[]") {
		// Make pointer for optional non-array types
		if !strings.Contains(baseType, "*") {
			baseType = "*" + baseType
		}
	}

	// Handle nillable
	if elem.Nillable && !strings.HasPrefix(baseType, "*") && !strings.HasPrefix(baseType, "[]") {
		baseType = "*" + baseType
	}

	return baseType
}

// buildXMLTag builds the XML tag for an element
func (ctg *ComplexTypeGenerator) buildXMLTag(elem models.Element) string {
	tag := elem.Name

	// Add omitempty for optional elements
	if elem.MinOccurs == "0" {
		tag += ",omitempty"
	}

	return tag
}

// IsComplexType checks if a type is a complex type (not a primitive)
func IsComplexType(typeName string) bool {
	primitives := map[string]bool{
		"string": true, "int": true, "integer": true, "long": true,
		"short": true, "byte": true, "boolean": true, "bool": true,
		"float": true, "double": true, "decimal": true,
		"dateTime": true, "date": true, "time": true,
		"base64Binary": true, "hexBinary": true,
	}

	// Remove namespace prefix
	if idx := strings.LastIndex(typeName, ":"); idx != -1 {
		typeName = typeName[idx+1:]
	}

	return !primitives[typeName]
}

// GenerateArrayHelper generates helper functions for array handling
func (ctg *ComplexTypeGenerator) GenerateArrayHelper(typeName string) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("// %sArray is a helper type for array marshaling\n", typeName))
	b.WriteString(fmt.Sprintf("type %sArray []%s\n\n", typeName, typeName))

	b.WriteString(fmt.Sprintf("// MarshalXML implements xml.Marshaler for %sArray\n", typeName))
	b.WriteString(fmt.Sprintf("func (a %sArray) MarshalXML(e *xml.Encoder, start xml.StartElement) error {\n", typeName))
	b.WriteString("\tfor _, item := range a {\n")
	b.WriteString("\t\tif err := e.EncodeElement(item, start); err != nil {\n")
	b.WriteString("\t\t\treturn err\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t}\n")
	b.WriteString("\treturn nil\n")
	b.WriteString("}\n\n")

	return b.String()
}
