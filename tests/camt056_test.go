package tests

import (
	"encoding/xml"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/mbanq/iso20022-go/pkg/common"
	"github.com/mbanq/iso20022-go/pkg/fednow"
	"github.com/mbanq/iso20022-go/pkg/fednow/bah"
	"github.com/mbanq/iso20022-go/pkg/fednow/camt"
	"github.com/mbanq/iso20022-go/pkg/fednow/config"
)

func TestCamt056_ParseViaFednowParse(t *testing.T) {
	cfg, err := config.LoadConfig("../config.json")
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	// Build a minimal camt.056 JSON struct, then construct XML envelope similar to generator output.
	msg := camt.FedNowMessageCxlReq{
		FedNowMsg: camt.FedNowCxlReq{
			CreationDateTime: common.ISODateTime(time.Now()),
			Identifier: camt.FedNowIdentifier{
				BusinessMessageID: "BizMsgId-TEST-CAMT056",
				MessageID:         "MsgId-TEST-CAMT056",
				MessageType:       "camt.056.001.08",
			},
			OriginalIdentifier: camt.FedNowIdentifier{
				MessageID:   "OrigMsgId-TEST",
				MessageType: "pacs.008.001.08",
				EndToEndID:  "E2E-TEST",
			},
			SenderDI: camt.FedNowDepositoryInstitution{
				SenderABANumber: "121182904",
			},
			ReceiverDI: camt.FedNowDepositoryInstitution{
				ReceiverABANumber: "084106768",
			},
		},
	}

	appHdr, err := bah.BuildBah(string(msg.FedNowMsg.Identifier.MessageID), cfg, "camt.056.001.08")
	if err != nil {
		t.Fatalf("failed to build AppHdr: %v", err)
	}

	document, err := camt.BuildCamt056Struct(msg, cfg)
	if err != nil {
		t.Fatalf("failed to build camt.056 document: %v", err)
	}

	appHdrPayload, err := xml.MarshalIndent(appHdr, "  ", "  ")
	if err != nil {
		t.Fatalf("failed to marshal AppHdr: %v", err)
	}
	appHdrXML := strings.Replace(string(appHdrPayload), "<BusinessApplicationHeaderV02>", "<AppHdr xmlns=\"urn:iso:std:iso:20022:tech:xsd:head.001.001.02\">", 1)
	appHdrXML = strings.Replace(appHdrXML, "</BusinessApplicationHeaderV02>", "</AppHdr>", 1)

	docPayload, err := xml.MarshalIndent(document, "  ", "  ")
	if err != nil {
		t.Fatalf("failed to marshal camt.056 document: %v", err)
	}
	docXML := string(docPayload)
	if strings.Contains(docXML, "<Document>") {
		docXML = strings.Replace(docXML, "<Document>", "<Document xmlns=\"urn:iso:std:iso:20022:tech:xsd:camt.056.001.08\">", 1)
	}

	// Use a simple root wrapper. fednow.Parse only needs AppHdr first, then Document.
	envelope := fmt.Sprintf("<Envelope>\n%s\n%s\n</Envelope>", appHdrXML, docXML)

	parsed, err := fednow.Parse([]byte(envelope))
	if err != nil {
		t.Fatalf("fednow.Parse failed: %v", err)
	}

	// Ensure the parsed type is the camt.056 wrapper we expect.
	if _, ok := parsed.(*camt.FedNowMessageCxlReq); !ok {
		t.Fatalf("expected camt.FedNowMessageCxlReq, got %T", parsed)
	}

}

