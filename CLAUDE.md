# ISO 20022 Go Library

## Overview

Go library for building, parsing, and validating ISO 20022 financial messages, with FedNow-specific extensions. Used as a dependency by `fednow-service-go`.

- **Module**: `github.com/mbanq/iso20022-go`
- **Go version**: 1.22
- **Type**: Library (with CLI tools)

## Package Structure

```
pkg/
  fednow/              # FedNow-specific message handling (main package)
    fednowmessage.go   # Main message interface
    generator.go       # Message generation
    parser.go          # Message parsing
    pacs/              # PACS message handlers (payments clearing & settlement)
    admi/              # ADMI message handlers (administrative)
    camt/              # CAMT message handlers (cash management)
    pain/              # PAIN message handlers (payments initiation)
    bah/               # Business Application Header handling
    config/            # Configuration management
  common/              # Common types and utilities (DateTime, etc.)

ISO20022/              # Generated Go models from XSD schemas (models.go per message type)
Internal/XSD/          # XML Schema Definition files
  iso/                 # Standard ISO XSD schemas
  fednow/              # FedNow-specific XSD schemas

cmd/                   # CLI tools
  converter/           # XML/message format converter
  generator/           # Message generation tool
  parser/              # Message parsing tool
  keys/                # Key management tool
  pacs008demo/         # PACS.008 demo
  pacs004demo/         # PACS.004 demo

scripts/               # Build and code generation scripts
  generate.sh          # XSD to Go model generation
  fix_imports.go       # Post-generation import fixes
  fix_inner_xml.go     # Post-generation XML handling fixes

sample_files/          # Sample XML and JSON message examples
tests/                 # Unit tests
```

## Supported ISO 20022 Message Types

### PACS (Payments Clearing & Settlement)
- **pacs.002.001.10** - Payment Status Report
- **pacs.004.001.10** - Payment Return
- **pacs.008.001.08** - Customer Credit Transfer
- **pacs.009.001.08** - Financial Institution Credit Transfer
- **pacs.028.001.03** - Payment Status Request

### CAMT (Cash Management)
- **camt.026.001.07** - Unable to Apply
- **camt.028.001.09** - Additional Payment Information
- **camt.029.001.09** - Resolution of Investigation
- **camt.052.001.08** - Bank to Customer Account Report
- **camt.054.001.08** - Bank to Customer Debit/Credit Notification
- **camt.055.001.09** - Customer Payment Cancellation Request
- **camt.056.001.08** - FI to FI Payment Cancellation Request
- **camt.060.001.05** - Account Reporting Request

### ADMI (Administration)
- **admi.002.001.01** - Message Reject
- **admi.004.001.02** - System Event Notification
- **admi.006.001.01** - Resend Request
- **admi.007.001.01** - Receipt Acknowledgement
- **admi.011.001.01** - System Event Acknowledgement
- **admi.998.001.02** - FedNow Participant File (FedNow-specific)

### PAIN (Payments Initiation)
- **pain.013.001.07** - Creditor Payment Activation Request
- **pain.014.001.07** - Creditor Payment Activation Request Status

### HEAD (Header)
- **head.001.001.02** - Business Application Header

## Key Interfaces

- `fednow.Generator` - Generate ISO 20022 XML messages
- `fednow.Parser` - Parse ISO 20022 XML into Go structs
- `fednow.FedNowMessage` - Main message interface

## Code Generation

Models in `ISO20022/` are generated from XSD schemas:

```bash
./scripts/generate.sh
```

Do NOT manually edit files in `ISO20022/` - they are auto-generated.

## Testing

```bash
go test ./tests/...
```

## How It's Consumed

Imported by `fednow-service-go` as:
```go
import "github.com/mbanq/iso20022-go"
```

Key sub-packages used:
- `pacs_008_001_08` - Credit transfer initiation
- `pacs_002_001_10` - Payment status
- `head_001_001_02` - Application header
- `admi_002_001_01`, `admi_007_001_01` - Admin messages
- `camt_029_001_09`, `camt_056_001_08` - Cancellation messages
- `pkg/fednow` - FedNow builders and converters
- `pkg/common` - Common ISO 20022 types

## Related Projects

- **fednow-service-go** (`/Users/vamsigv/Documents/Mbanq/fednow-service-go`): Main FedNow payment processing service
- **mq-client** (`/Users/vamsigv/Documents/Mbanq/mq-client`): MQ messaging gateway
