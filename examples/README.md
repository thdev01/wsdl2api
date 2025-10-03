## Examples

This directory contains example WSDLs for testing wsdl2api.

### Correios (Brazilian Postal Service)

```bash
# Calculate shipping costs
wsdl2api serve --wsdl https://apps.correios.com.br/SigepMasterJPA/AtendeClienteService/AtendeCliente?wsdl --port 8080
```

### ViaCEP (Brazilian Postal Code Lookup)

```bash
# CEP lookup service
wsdl2api serve --wsdl http://viacep.com.br/ws/{cep}/xml/ --port 8080
```

### Public SOAP Services for Testing

1. **Currency Converter**
   - WSDL: http://www.webservicex.net/CurrencyConvertor.asmx?WSDL
   - Description: Convert between currencies

2. **Weather Service**
   - WSDL: http://www.webservicex.net/globalweather.asmx?WSDL
   - Description: Global weather information

3. **Calculator**
   - WSDL: http://www.dneonline.com/calculator.asmx?WSDL
   - Description: Simple calculator operations

### Usage Examples

#### Generate Code

```bash
wsdl2api generate \
  --wsdl http://www.dneonline.com/calculator.asmx?WSDL \
  --output ./generated/calculator \
  --package calculator
```

#### Start Server

```bash
wsdl2api serve \
  --wsdl http://www.dneonline.com/calculator.asmx?WSDL \
  --port 8080
```

#### Test Endpoint

```bash
# Get service info
curl http://localhost:8080/info

# Execute operation
curl -X POST http://localhost:8080/api/Add \
  -H "Content-Type: application/json" \
  -d '{"intA": 5, "intB": 3}'
```
