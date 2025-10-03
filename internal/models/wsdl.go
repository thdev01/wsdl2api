package models

// Definitions represents a WSDL definitions structure
type Definitions struct {
	Name            string
	TargetNamespace string
	Services        []Service
	Bindings        []Binding
	PortTypes       []PortType
	Messages        []Message
	Types           []Type
}

// Service represents a WSDL service
type Service struct {
	Name  string
	Ports []Port
}

// Port represents a service port
type Port struct {
	Name    string
	Binding string
	Address string
}

// Binding represents a WSDL binding
type Binding struct {
	Name       string
	Type       string
	Operations []BindingOperation
}

// BindingOperation represents an operation in a binding
type BindingOperation struct {
	Name      string
	SoapAction string
	Input     BindingMessage
	Output    BindingMessage
}

// BindingMessage represents input/output binding
type BindingMessage struct {
	Use       string
	Namespace string
}

// PortType represents a WSDL port type
type PortType struct {
	Name       string
	Operations []Operation
}

// Operation represents a WSDL operation
type Operation struct {
	Name          string
	Documentation string
	Input         Message
	Output        Message
}

// Message represents a WSDL message
type Message struct {
	Name  string
	Parts []Part
}

// Part represents a message part
type Part struct {
	Name    string
	Element string
	Type    string
}

// Type represents a WSDL/XSD type
type Type struct {
	Name       string
	Elements   []Element
	Attributes []Attribute
}

// Element represents an XSD element
type Element struct {
	Name      string
	Type      string
	MinOccurs string
	MaxOccurs string
	Nillable  bool
}

// Attribute represents an XSD attribute
type Attribute struct {
	Name string
	Type string
	Use  string
}
