package fednow

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"strings"

	admi002 "github.com/mbanq/iso20022-go/ISO20022/admi_002_001_01"
	head "github.com/mbanq/iso20022-go/ISO20022/head_001_001_02"
	pacs002 "github.com/mbanq/iso20022-go/ISO20022/pacs_002_001_10"
	pacs008 "github.com/mbanq/iso20022-go/ISO20022/pacs_008_001_08"
	"github.com/mbanq/iso20022-go/pkg/fednow/admi"
	"github.com/mbanq/iso20022-go/pkg/fednow/pacs"
)

// Parse an incoming pacs.008 XML file and return a JSON representation.
func Parse(xmlData []byte) ([]byte, error) {
	decoder := xml.NewDecoder(bytes.NewReader(xmlData))
	var appHdr head.BusinessApplicationHeaderV02
	foundAppHdr := false

	// First, find and decode the AppHdr
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if se, ok := token.(xml.StartElement); ok && se.Name.Local == "AppHdr" {
			if err := decoder.DecodeElement(&appHdr, &se); err != nil {
				return nil, err
			}
			foundAppHdr = true
			break // Stop after finding AppHdr
		}
	}

	if !foundAppHdr {
		return nil, errors.New("failed to find AppHdr in XML")
	}

	// Now, decode the Document based on the message type from AppHdr
	var fednowMsg interface{}
	var err error

	msgType := string(appHdr.MsgDefIdr)

	switch {
	case strings.Contains(msgType, "pacs.008.001.08"):
		var doc pacs008.Document
		if err = decoder.Decode(&doc); err != nil {
			return nil, err
		}
		fednowMsg, err = pacs.ParsePacs008(appHdr, doc)
	case strings.Contains(msgType, "pacs.002.001.10"):
		var doc pacs002.Document
		if err = decoder.Decode(&doc); err != nil {
			return nil, err
		}
		fednowMsg, err = pacs.ParsePacs002(appHdr, doc)
	case strings.Contains(msgType, "admi.002.001.01"):
		var doc admi002.Document
		if err = decoder.Decode(&doc); err != nil {
			return nil, err
		}
		fednowMsg, err = admi.ParseAdmi002Struct(&doc, appHdr)
	default:
		return nil, errors.New("unsupported message type: " + msgType)
	}

	if err != nil {
		return nil, err
	}

	// Marshal the struct to JSON
	return json.MarshalIndent(fednowMsg, "", "  ")
}
