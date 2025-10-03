# Public WSDL Catalog

A curated list of free, working public WSDL services for testing wsdl2api.

---

## ✅ Tested & Working

### 1. Calculator Service
- **URL**: http://www.dneonline.com/calculator.asmx?WSDL
- **Description**: Simple calculator with Add, Subtract, Multiply, Divide
- **Best For**: Beginners, basic testing
- **Operations**: 4

```bash
./wsdl2api serve --wsdl http://www.dneonline.com/calculator.asmx?WSDL --port 8080
```

### 2. Country Info Service
- **URL**: http://webservices.oorsprong.org/websamples.countryinfo/CountryInfoService.wso?WSDL
- **Description**: Country information (capital, currency, flag, phone codes)
- **Best For**: Complex data structures, multiple operations
- **Operations**: 20+

```bash
./wsdl2api serve --wsdl http://webservices.oorsprong.org/websamples.countryinfo/CountryInfoService.wso?WSDL --port 8080
```

### 3. Number Conversion Service
- **URL**: https://www.dataaccess.com/webservicesserver/numberconversion.wso?WSDL
- **Description**: Convert numbers to words and dollar amounts
- **Best For**: String manipulation, data transformation
- **Operations**: 2

```bash
./wsdl2api serve --wsdl https://www.dataaccess.com/webservicesserver/numberconversion.wso?WSDL --port 8080
```

### 4. Temperature Conversion Service
- **URL**: http://webservices.daehosting.com/services/TemperatureConversions.wso?WSDL
- **Description**: Convert between Celsius, Fahrenheit, Kelvin
- **Best For**: Unit conversion
- **Operations**: 6

```bash
./wsdl2api serve --wsdl http://webservices.daehosting.com/services/TemperatureConversions.wso?WSDL --port 8080
```

### 5. Hello World Service
- **URL**: http://www.learnwebservices.com/services/hello?WSDL
- **Description**: Simple greeting service
- **Best For**: Minimal testing
- **Operations**: 1

```bash
./wsdl2api serve --wsdl http://www.learnwebservices.com/services/hello?WSDL --port 8080
```

### 6. Bank BLZ Service (Germany)
- **URL**: http://www.thomas-bayer.com/axis2/services/BLZService?wsdl
- **Description**: German bank code lookup
- **Best For**: Real-world data lookup
- **Operations**: 2

```bash
./wsdl2api serve --wsdl http://www.thomas-bayer.com/axis2/services/BLZService?wsdl --port 8080
```

---

## 🇧🇷 Brazilian Services

### 7. Correios (Brazilian Postal Service)
- **URL**: https://apps.correios.com.br/SigepMasterJPA/AtendeClienteService/AtendeCliente?wsdl
- **Description**: CEP lookup, shipping calculation, tracking
- **Best For**: Complex enterprise SOAP services
- **Operations**: 30+
- **Note**: Some operations require authentication

```bash
./wsdl2api serve --wsdl https://apps.correios.com.br/SigepMasterJPA/AtendeClienteService/AtendeCliente?wsdl --port 8080
```

#### Free Operations:
- `consultaCEP` - CEP (postal code) lookup
- `getStatusCartaoPostagem` - Check posting card status

#### Alternative (REST):
- **ViaCEP**: https://viacep.com.br/ (No WSDL, REST only)
- **BrasilAPI**: https://brasilapi.com.br/ (No WSDL, REST only)

---

## ⚠️ Currently Unavailable

### WebServiceX Services (Offline)
Many popular WebServiceX.net services are no longer operational:
- ❌ Global Weather: http://www.webservicex.net/globalweather.asmx?WSDL
- ❌ Currency Converter: http://www.webservicex.net/CurrencyConvertor.asmx?WSDL
- ❌ Stock Quote: http://www.webservicex.net/stockquote.asmx?WSDL

---

## 🧪 Testing Strategy

### Quick Test (Calculator)
```bash
./wsdl2api serve --wsdl examples/calculator.wsdl --port 8080
curl http://localhost:8080/info
```

### Complex Test (Country Info)
```bash
./wsdl2api serve --wsdl http://webservices.oorsprong.org/websamples.countryinfo/CountryInfoService.wso?WSDL --port 8080
curl http://localhost:8080/info
```

### Real-World Test (Correios)
```bash
./wsdl2api serve --wsdl https://apps.correios.com.br/SigepMasterJPA/AtendeClienteService/AtendeCliente?wsdl --port 8080
curl http://localhost:8080/info
```

---

## 📊 WSDL Complexity Levels

| Service | Complexity | Operations | Types | Best For |
|---------|-----------|------------|-------|----------|
| Hello World | ⭐ Basic | 1 | 2 | Getting started |
| Calculator | ⭐ Basic | 4 | 4 | Simple testing |
| Number Conversion | ⭐⭐ Medium | 2 | 4 | Data transformation |
| Temperature | ⭐⭐ Medium | 6 | 6 | Multiple operations |
| Bank BLZ | ⭐⭐⭐ Advanced | 2 | 10+ | Complex types |
| Country Info | ⭐⭐⭐ Advanced | 20+ | 30+ | Large service |
| Correios | ⭐⭐⭐⭐ Expert | 30+ | 50+ | Enterprise SOAP |

---

## 🔍 Finding More WSDLs

### Search Tips:
1. Google: `site:*.gov.br filetype:wsdl`
2. GitHub: `filename:*.wsdl`
3. SOAP Testing Tools: Check SoapUI example projects

### Validation:
```bash
# Test if WSDL is accessible
curl -I <wsdl-url>

# Parse with wsdl2api
./wsdl2api generate --wsdl <wsdl-url> --output ./test
```

---

**Last Updated**: 2025-10-03
