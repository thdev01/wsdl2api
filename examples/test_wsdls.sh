#!/bin/bash

# Test script for various WSDL examples
# Tests that wsdl2api can parse and serve different WSDLs

set -e

WSDL2API="../wsdl2api"
PORT=8080

echo "Testing WSDL2API with various examples..."
echo ""

# Test 1: Local Calculator WSDL
echo "Test 1: Calculator (local file)"
timeout 3s $WSDL2API serve --wsdl calculator.wsdl --port $PORT 2>&1 | grep -q "Found 1 services" && echo "✅ PASS" || echo "❌ FAIL"
echo ""

# Test 2: Number Conversion WSDL
echo "Test 2: Number Conversion (local file)"
timeout 3s $WSDL2API serve --wsdl numberconversion.wsdl --port $PORT 2>&1 | grep -q "Found 1 services" && echo "✅ PASS" || echo "❌ FAIL"
echo ""

# Test 3: Temperature WSDL
echo "Test 3: Temperature Conversion (local file)"
timeout 3s $WSDL2API serve --wsdl temperature.wsdl --port $PORT 2>&1 | grep -q "Found 1 services" && echo "✅ PASS" || echo "❌ FAIL"
echo ""

# Test 4: Country Info (remote WSDL)
echo "Test 4: Country Info Service (remote URL)"
timeout 5s $WSDL2API serve --wsdl "http://webservices.oorsprong.org/websamples.countryinfo/CountryInfoService.wso?WSDL" --port $PORT 2>&1 | grep -q "Found 1 services" && echo "✅ PASS" || echo "❌ FAIL"
echo ""

# Test 5: Code generation
echo "Test 5: Code Generation (calculator)"
$WSDL2API generate --wsdl calculator.wsdl --output /tmp/test-gen --package testclient > /dev/null 2>&1 && \
  [ -f /tmp/test-gen/client.go ] && \
  echo "✅ PASS - Generated client.go" || echo "❌ FAIL"
rm -rf /tmp/test-gen
echo ""

echo "All tests completed!"
