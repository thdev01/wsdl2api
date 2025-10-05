package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/thdev01/wsdl2api/pkg/exporter"
	"github.com/thdev01/wsdl2api/pkg/generator"
	"github.com/thdev01/wsdl2api/pkg/parser"
	"github.com/thdev01/wsdl2api/pkg/server"
	"github.com/thdev01/wsdl2api/pkg/typescript"
)

var (
	wsdlPath         string
	outputDir        string
	packageName      string
	port             int
	host             string
	exportFormat     string
	generateMock     bool
	soapVersion      string
	generateTS       bool
	tsOutputDir      string
)

var rootCmd = &cobra.Command{
	Use:   "wsdl2api",
	Short: "Convert WSDL to REST API",
	Long:  `WSDL2API converts legacy SOAP/WSDL services into modern REST APIs`,
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Go client code from WSDL",
	Long:  `Parse WSDL and generate complete Go client structures`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if wsdlPath == "" {
			return fmt.Errorf("wsdl path is required")
		}

		fmt.Printf("Parsing WSDL: %s\n", wsdlPath)

		// Parse WSDL
		p := parser.NewParser()
		definitions, err := p.Parse(wsdlPath)
		if err != nil {
			return fmt.Errorf("failed to parse WSDL: %w", err)
		}

		fmt.Printf("Found %d services\n", len(definitions.Services))

		// Generate code
		g := generator.NewGenerator(outputDir, packageName)
		if generateMock {
			if err := g.GenerateWithMock(definitions); err != nil {
				return fmt.Errorf("failed to generate code: %w", err)
			}
		} else {
			if err := g.Generate(definitions); err != nil {
				return fmt.Errorf("failed to generate code: %w", err)
			}
		}

		fmt.Printf("Code generated successfully in: %s\n", outputDir)
		return nil
	},
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start REST API server",
	Long:  `Parse WSDL, generate code, and start REST API server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if wsdlPath == "" {
			return fmt.Errorf("wsdl path is required")
		}

		fmt.Printf("Parsing WSDL: %s\n", wsdlPath)

		// Parse WSDL
		p := parser.NewParser()
		definitions, err := p.Parse(wsdlPath)
		if err != nil {
			return fmt.Errorf("failed to parse WSDL: %w", err)
		}

		fmt.Printf("Found %d services\n", len(definitions.Services))

		// Start server
		srv := server.NewServer(definitions, host, port)
		fmt.Printf("Starting REST API server on %s:%d\n", host, port)

		if err := srv.Start(); err != nil {
			return fmt.Errorf("failed to start server: %w", err)
		}

		return nil
	},
}

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export WSDL to OpenAPI/Swagger specification",
	Long:  `Parse WSDL and export as OpenAPI 3.0 specification`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if wsdlPath == "" {
			return fmt.Errorf("wsdl path is required")
		}

		fmt.Printf("Parsing WSDL: %s\n", wsdlPath)

		// Parse WSDL
		p := parser.NewParser()
		definitions, err := p.Parse(wsdlPath)
		if err != nil {
			return fmt.Errorf("failed to parse WSDL: %w", err)
		}

		fmt.Printf("Converting to OpenAPI...\n")

		// Convert to OpenAPI
		spec, err := exporter.ConvertWSDLToOpenAPI(definitions)
		if err != nil {
			return fmt.Errorf("failed to convert to OpenAPI: %w", err)
		}

		// Export based on format
		var output string
		if exportFormat == "yaml" || exportFormat == "yml" {
			output, err = spec.ExportToYAML()
		} else {
			output, err = spec.ExportToJSON()
		}

		if err != nil {
			return fmt.Errorf("failed to export: %w", err)
		}

		// Write to file or stdout
		if outputDir == "" || outputDir == "-" {
			fmt.Println(output)
		} else {
			filename := fmt.Sprintf("%s/openapi.%s", outputDir, exportFormat)
			if err := os.WriteFile(filename, []byte(output), 0644); err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}
			fmt.Printf("OpenAPI spec exported to: %s\n", filename)
		}

		// Generate TypeScript client if requested
		if generateTS {
			tsDir := tsOutputDir
			if tsDir == "" {
				tsDir = filepath.Join(outputDir, "typescript")
			}

			fmt.Printf("Generating TypeScript client in: %s\n", tsDir)
			tsGen := typescript.NewGenerator(tsDir, spec)
			if err := tsGen.Generate(); err != nil {
				return fmt.Errorf("failed to generate TypeScript client: %w", err)
			}
			fmt.Printf("TypeScript client generated successfully in: %s\n", tsDir)
		}

		return nil
	},
}

func init() {
	// Generate command flags
	generateCmd.Flags().StringVarP(&wsdlPath, "wsdl", "w", "", "WSDL file path or URL (required)")
	generateCmd.Flags().StringVarP(&outputDir, "output", "o", "./generated", "Output directory")
	generateCmd.Flags().StringVarP(&packageName, "package", "p", "client", "Go package name")
	generateCmd.Flags().BoolVar(&generateMock, "mock", false, "Generate mock server")
	generateCmd.Flags().StringVar(&soapVersion, "soap-version", "1.1", "SOAP version (1.1 or 1.2)")
	_ = generateCmd.MarkFlagRequired("wsdl")

	// Serve command flags
	serveCmd.Flags().StringVarP(&wsdlPath, "wsdl", "w", "", "WSDL file path or URL (required)")
	serveCmd.Flags().IntVar(&port, "port", 8080, "Server port")
	serveCmd.Flags().StringVar(&host, "host", "localhost", "Server host")
	_ = serveCmd.MarkFlagRequired("wsdl")

	// Export command flags
	exportCmd.Flags().StringVarP(&wsdlPath, "wsdl", "w", "", "WSDL file path or URL (required)")
	exportCmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory (empty for stdout)")
	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "json", "Export format (json or yaml)")
	exportCmd.Flags().BoolVar(&generateTS, "typescript", false, "Generate TypeScript client")
	exportCmd.Flags().StringVar(&tsOutputDir, "ts-output", "", "TypeScript output directory (default: <output>/typescript)")
	_ = exportCmd.MarkFlagRequired("wsdl")

	// Add commands to root
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(exportCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
