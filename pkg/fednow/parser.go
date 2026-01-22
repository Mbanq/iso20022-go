package fednow

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"strings"

	admi002 "github.com/mbanq/iso20022-go/ISO20022/admi_002_001_01"
	admi007 "github.com/mbanq/iso20022-go/ISO20022/admi_007_001_01"
	camt056 "github.com/mbanq/iso20022-go/ISO20022/camt_056_001_08"
	head "github.com/mbanq/iso20022-go/ISO20022/head_001_001_02"
	pacs002 "github.com/mbanq/iso20022-go/ISO20022/pacs_002_001_10"
	pacs008 "github.com/mbanq/iso20022-go/ISO20022/pacs_008_001_08"
	pain013 "github.com/mbanq/iso20022-go/ISO20022/pain_013_001_07"
	"github.com/mbanq/iso20022-go/pkg/fednow/admi"
	"github.com/mbanq/iso20022-go/pkg/fednow/camt"
	"github.com/mbanq/iso20022-go/pkg/fednow/pacs"
	"github.com/mbanq/iso20022-go/pkg/fednow/pain"
)

// Parse an incoming pacs.008 XML file and return a JSON representation.
func Parse(xmlData []byte) (FedNowMessage, error) {
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
	var fednowMsg FedNowMessage
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
	case strings.Contains(msgType, "admi.007.001.01"):
		var doc admi007.Document
		if err = decoder.Decode(&doc); err != nil {
			return nil, err
		}
		fednowMsg, err = admi.ParseAdmi007Struct(&doc, appHdr)
	case strings.Contains(msgType, "pain.013.001.07"):
		var doc pain013.Document
		if err = decoder.Decode(&doc); err != nil {
			return nil, err
		}
		fednowMsg, err = pain.ParsePain013(appHdr, doc)
	case strings.Contains(msgType, "camt.056.001.08"):
		var doc camt056.Document
		if err = decoder.Decode(&doc); err != nil {
			return nil, err
		}
		fednowMsg, err = camt.ParseCamt056(appHdr, doc)
	default:
		return nil, errors.New("unsupported message type: " + msgType)
	}

	if err != nil {
		return nil, err
	}

	// Marshal the struct to JSON
	return fednowMsg, nil
}
