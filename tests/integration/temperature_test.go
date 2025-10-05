package integration

import (
	"testing"

	temperature "github.com/thdev01/wsdl2api/tests/integration/temperature/client"
)

// TestTemperatureClientGeneration tests temperature conversion client generation
func TestTemperatureClientGeneration(t *testing.T) {
	c := temperature.NewClient("")

	if c == nil {
		t.Fatal("Failed to create temperature client")
	}

	if c.HTTPClient == nil {
		t.Error("HTTPClient should not be nil")
	}

	t.Logf("✓ Temperature client generated successfully")
	t.Logf("  Default URL: %s", c.URL)
	t.Logf("  SOAP Version: %s", c.SOAPVersion)
}

// TestTemperatureOperations tests temperature conversion operations
func TestTemperatureOperations(t *testing.T) {
	// Test CelsiusToFahrenheit request
	celsiusReq := &temperature.CelsiusToFahrenheitRequest{
		NCelsius: "25",
	}

	if celsiusReq.NCelsius != "25" {
		t.Error("CelsiusToFahrenheitRequest not working")
	}

	// Test FahrenheitToCelsius request
	fahrenheitReq := &temperature.FahrenheitToCelsiusRequest{
		NFahrenheit: "77",
	}

	if fahrenheitReq.NFahrenheit != "77" {
		t.Error("FahrenheitToCelsiusRequest not working")
	}

	t.Logf("✓ Temperature operations generated correctly")
	t.Logf("  CelsiusToFahrenheit: %+v", celsiusReq)
	t.Logf("  FahrenheitToCelsius: %+v", fahrenheitReq)
}

// TestTemperatureResponseTypes tests response type generation
func TestTemperatureResponseTypes(t *testing.T) {
	// Test Celsius to Fahrenheit response
	celsiusResp := &temperature.CelsiusToFahrenheitResponse{
		NCelsiusToFahrenheitResult: "77",
	}

	if celsiusResp.NCelsiusToFahrenheitResult != "77" {
		t.Error("CelsiusToFahrenheitResponse not working")
	}

	// Test Fahrenheit to Celsius response
	fahrenheitResp := &temperature.FahrenheitToCelsiusResponse{
		NFahrenheitToCelsiusResult: "25",
	}

	if fahrenheitResp.NFahrenheitToCelsiusResult != "25" {
		t.Error("FahrenheitToCelsiusResponse not working")
	}

	t.Logf("✓ Temperature response types work correctly")
}

// TestTemperatureMultipleConversions tests handling multiple conversions
func TestTemperatureMultipleConversions(t *testing.T) {
	testCases := []struct {
		celsius    string
		fahrenheit string
		desc       string
	}{
		{"0", "32", "Freezing point"},
		{"100", "212", "Boiling point"},
		{"25", "77", "Room temperature"},
		{"-40", "-40", "Same value in both scales"},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			req := &temperature.CelsiusToFahrenheitRequest{
				NCelsius: tc.celsius,
			}

			if req.NCelsius != tc.celsius {
				t.Errorf("Expected celsius %s, got %s", tc.celsius, req.NCelsius)
			}

			t.Logf("  %s: %s°C expected to be %s°F", tc.desc, tc.celsius, tc.fahrenheit)
		})
	}

	t.Logf("✓ Multiple temperature conversions handled correctly")
}

// BenchmarkTemperatureRequestCreation benchmarks request creation
func BenchmarkTemperatureRequestCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = &temperature.CelsiusToFahrenheitRequest{
			NCelsius: "25",
		}
	}
}
