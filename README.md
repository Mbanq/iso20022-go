# ISO20022-GO Library

A comprehensive ISO20022 library built specifically for US Payment Rails, with **FedNow** as the first supported payment rail. This Go library provides robust message conversion capabilities with custom JSON payloads for seamless integration with modern payment systems.

## Features

- **Complete ISO20022 Message Support**: Supports PACS message conversion with custom JSON payloads
- **FedNow Integration**: Built-in support for FedNow payment rail with envelope wrapping
- **Bi-directional Conversion**: Convert between JSON and XML formats with proper namespace handling
- **Transmission-Ready XML**: Generate XML messages with FedNow envelopes ready for transmission
- **Flexible Parser**: Parse incoming XML messages to custom JSON payloads
- **XSD-Based Validation**: Dynamic XSD parsing for accurate message structure validation

## Project Structure

```
iso20022-go/
├── cmd/                              # Command-line utilities and demos
│   ├── converter/                    # JSON to XML converter (without FedNow envelope)
│   ├── generator/                    # Demo for generating transmission-ready XML
│   ├── parser/                       # Demo for parsing XML to JSON
│   ├── pacs008demo/                  # PACS 008 message demo
│   └── pacs004demo/                  # PACS 004 message demo
├── ISO20022/                         # Generated ISO20022 message structures
│   ├── pacs_008_001_08/              # PACS 008 Credit Transfer message models
│   ├── pacs_002_001_10/              # PACS 002 Payment Status Report models
│   ├── pacs_004_001_10/              # PACS 004 Payment Return models
│   ├── head_001_001_02/              # Business Application Header models
│   ├── camt_*/                       # Cash Management message models
│   ├── pain_*/                       # Payment Initiation message models
│   └── admi_*/                       # Administrative message models
├── pkg/                              # Core library packages
│   ├── fednow/                       # FedNow-specific functionality
│   │   ├── generator.go              # Transmission-ready XML generation
│   │   ├── parser.go                 # XML to JSON parsing
│   │   ├── pacs/                     # PACS message builders
│   │   ├── bah/                      # Business Application Header builders
│   │   └── config/                   # Configuration structures
│   └── common/                       # Shared utilities and helpers
├── Internal/                         # Internal XSD files and schemas
│   └── XSD/                          # XSD schema files for validation
├── sample_files/                     # Sample JSON and XML files for testing
├── scripts/                          # Utility scripts for code generation
└── tests/                            # Test files and test utilities
```

## Installation

```bash
go get github.com/Mbanq/iso20022-go
```

## Usage

### 1. Generating Transmission-Ready XML Messages

Use the `Generate` function from `pkg/fednow/generator.go` to create FedNow-compliant XML messages with proper envelopes:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/Mbanq/iso20022-go/pkg/fednow"
    "github.com/Mbanq/iso20022-go/pkg/fednow/config"
    "github.com/Mbanq/iso20022-go/pkg/fednow/pacs"
)

func main() {
    // Create configuration
    cfg := &config.Config{
        // Add your FedNow configuration here
        ParticipantID: "YOUR_PARTICIPANT_ID",
        // ... other config fields
    }
    
    // Create your FedNow message (example with PACS 008)
    message := pacs.FedNowMessageCCT{
        // Populate your message fields
    }
    
    // Generate transmission-ready XML with FedNow envelope
    xmlData, err := fednow.Generate(
        "path/to/xsd/file.xsd",
        "pacs.008.001.08",
        cfg,
        message,
    )
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(string(xmlData))
}
```

**Supported Message Types:**
- `pacs.008.001.08` - Customer Credit Transfer
- `pacs.002.001.10` - Payment Status Report  
- `pacs.004.001.10` - Payment Return

### 2. Parsing XML Messages to Custom JSON

Use the `Parse` function from `pkg/fednow/parser.go` to convert incoming XML messages to custom JSON payloads:

```go
package main

import (
    "fmt"
    "io/ioutil"
    "log"
    
    "github.com/mbanq/iso20022-go/pkg/fednow"
)

func main() {
    // Read XML file
    xmlData, err := ioutil.ReadFile("incoming_message.xml")
    if err != nil {
        log.Fatal(err)
    }
    
    // Parse XML to custom JSON payload
    jsonData, err := fednow.Parse(xmlData)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(string(jsonData))
}
```

### 3. Simple JSON to XML Conversion (Without FedNow Envelope)

For basic ISO20022 message conversion without FedNow envelope wrapping, use the converter utility:

```bash
# Navigate to converter directory
cd cmd/converter

# Run the converter (modify the file paths in converter.go as needed)
go run converter.go
```

**Example converter usage:**

```go
package main

import (
    "encoding/json"
    "encoding/xml"
    "io/ioutil"
    "log"
    "strings"
    
    pacs "github.com/mbanq/iso20022-go/ISO20022/pacs_008_001_08"
)

func main() {
    // Read JSON file
    jsonFile, err := ioutil.ReadFile("sample.json")
    if err != nil {
        log.Fatal(err)
    }
    
    // Unmarshal JSON to Go struct
    var doc pacs.Document
    if err := json.Unmarshal(jsonFile, &doc); err != nil {
        log.Fatal(err)
    }
    
    // Marshal Go struct to XML
    xmlBytes, err := xml.MarshalIndent(doc, "", "  ")
    if err != nil {
        log.Fatal(err)
    }
    
    // Add proper namespace
    xmlString := strings.Replace(string(xmlBytes), 
        `<Document>`, 
        `<Document xmlns="urn:iso:std:iso:20022:tech:xsd:pacs.008.001.08">`, 
        1)
    
    // Add XML header and write to file
    output := []byte(xml.Header + xmlString)
    err = ioutil.WriteFile("output.xml", output, 0644)
    if err != nil {
        log.Fatal(err)
    }
}
```

## Running Examples

The library includes several demo applications:

```bash
# XML Parser Demo
cd cmd/parser && go run main.go

# XML Generator Demo  
cd cmd/generator && go run main.go
```

## Testing

```bash
# Run all tests
go test ./...

# Run specific test
go test ./tests/
```

## Configuration

The library uses a configuration structure for FedNow-specific settings. Create a `Config` struct with your participant details:

```go
type Config struct {
    ParticipantID     string
    BusinessDate      string
    SequenceNumber    string
    // ... additional configuration fields
}
```

## Message Flow

1. **Outbound Messages**: JSON → Go Struct → XML with FedNow Envelope → Transmission
2. **Inbound Messages**: XML with FedNow Envelope → Go Struct → Custom JSON Payload

## Supported ISO20022 Messages

- **PACS (Payment Clearing and Settlement)**
  - pacs.008.001.08 - Customer Credit Transfer - Available Now
  - pacs.002.001.10 - Payment Status Report - Available Now 
  - pacs.004.001.10 - Payment Return - Available Now
  - pacs.028.001.03 - FI To FI Payment Status Report - Available Now

- **CAMT (Cash Management)**
  - Multiple CAMT message types supported - WIP

- **PAIN (Payment Initiation)**
  - pain.013.001.07 - Customer Credit Transfer Initiation - WIP
  - pain.014.001.07 - Customer Credit Transfer Cancellation Request - WIP

- **ADMI (Administration)**
  - admi.002.001.01 - Payment Status Report(Parsing only) - Available Now
  - admi.004.001.01 - Payment Return(Parsing only) - WIP

## Requirements

- Go 1.22 or higher
- No external dependencies (uses only Go standard library)

## US Payment Rails

This project can be used to generate the message content for US Payment Rails, which utilize the ISO 20022 standard. However, the message envelope format for US Payment Rails is proprietary. Users must obtain the necessary proprietary XSD files from the Federal Reserve and integrate them to create complete US Payment Rails messages. This project does not include any proprietary information or files.


### Important Notice on Proprietary Envelope
The US Payment Rails envelope format is proprietary and must be implemented according to the Federal Reserve specifications. This library only generates the ISO 20022 message content and wraps them in the message envelope only if the propietory XSD is provided.

#### For Fednow:

[Get Access to Prorietory XSD & Fednow Documentation](https://www.frbservices.org/resources/financial-services/fednow/fednow-devrel)

For official information about Fednow ISO 20022 implementation, including envelope specifications, please consult the following resources:

- [To Learn More about Fednow](https://explore.fednow.org/guided-journey)
- [Fednow Documentation](https://www.frbservices.org/resources/financial-services/fednow/fednow-devrel)
- [My Standards](https://www2.swift.com/mystandards/#/group/Federal_Reserve_Financial_Services/FedNow_Service)


#### For Fedwire: 

[Get Access to Prorietory XSD](https://www.frbservices.org/binaries/content/assets/crsocms/resources/financial-services/request-access-fedwire-funds-iso-20022-technical-guide.pdf)

For official information about Fedwire ISO 20022 implementation, including envelope specifications, please consult the following resources:

- [Federal Reserve ISO 20022 Implementation Center](https://www.frbservices.org/resources/financial-services/wires/iso-20022-implementation-center)
- [MyStandards platform](https://www2.swift.com/mystandards/#/group/Federal_Reserve_Financial_Services/Fedwire_Funds_Service)
- [Fedwire Funds Service ISO 20022 Implementation Guide and Technical Guide](https://www.frbservices.org/resources/financial-services/wires/iso-20022-implementation-center)



## Acknowledgments

We would like to thank the following projects and resources that made this library possible:

- **ISO20022 Organization** - For providing comprehensive payment messaging standards. [ISO20022](https://www.iso20022.org/)
- **Federal Reserve** - For FedNow service specifications and documentation. [Federal Reserve](https://www.frbservices.org/)
- **MOOV** - For XSD2GO Module and fednow20022 Project which were crucial for building the Go structs from XSDs. [XSD2GO](https://github.com/moov-io/xsd2go) and [fednow20022](https://github.com/moov-io/fednow20022)
- **Go Community** - For excellent XML and JSON handling libraries
- **Contributors** - All developers who have contributed to this project

## Contributing

Contributions are welcome and will be accepting via pull requests soon. Please Star and Watch the Repository for updates. You can create the Pull Requests to raise any issues, vulnerabilities, suggestions or feature requests.

[GitHub Issues page](https://github.com/Mbanq/iso20022-go/issues).

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

---

**Note**: This library is specifically designed for US Payment Rails with FedNow as the primary target. For other payment rails or regions, additional configuration may be required.



<p align="center"><strong style="font-size: 2em">Built with ❤️ in the Beautiful State of Washington!</strong></p>