package tests

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/mbanq/iso20022-go/ISO20022/pacs_008_001_08"
	"github.com/mbanq/iso20022-go/pkg/common"
	"github.com/mbanq/iso20022-go/pkg/fednow"
	"github.com/mbanq/iso20022-go/pkg/fednow/config"
	"github.com/mbanq/iso20022-go/pkg/fednow/pacs"
)

func BenchmarkGenerate(b *testing.B) {
	// The path to the real XSD file, relative to the tests directory.
	xsdPath := "../Internal/XSD/fednow-outgoing_external.xsd"

	// Check if the XSD file exists to provide a better error message.
	if _, err := os.Stat(xsdPath); os.IsNotExist(err) {
		b.Fatalf("XSD file not found at %s. Make sure the path is correct.", xsdPath)
	}

	// Create mock config and message
	cfg, err := config.LoadConfig("../config.json")
	if err != nil {
		b.Fatalf("failed to load config: %v", err)
	}

	strPtr := func(s string) *pacs_008_001_08.Max35Text {
		t := pacs_008_001_08.Max35Text(s)
		return &t
	}

	strPtr140 := func(s string) *pacs_008_001_08.Max140Text {
		t := pacs_008_001_08.Max140Text(s)
		return &t
	}

	strPtr70 := func(s string) *pacs_008_001_08.Max70Text {
		t := pacs_008_001_08.Max70Text(s)
		return &t
	}

	strPtr16 := func(s string) *pacs_008_001_08.Max16Text {
		t := pacs_008_001_08.Max16Text(s)
		return &t
	}

	countryCodePtr := func(s string) *pacs_008_001_08.CountryCode {
		t := pacs_008_001_08.CountryCode(s)
		return &t
	}

	message := pacs.FedNowMessageCCT{
		FedNowMsg: pacs.FedNowDetails{
			CreationDateTime: common.ISODateTime(time.Now()),
			Identifier: pacs.FedNowIdentifier{
				BusinessMessageID: "20230604111111111Sc01Step1MsgId",
				MessageID:         "20230604111111111Sc01Step1MsgId",
				InstructionID:     strPtr("Scenario01InstrId001"),
				EndToEndID:        "Scenario01EtoEId001",
				TransactionID:     strPtr("BankARefNum000001"),
			},
			PaymentType: pacs.FedNowPaymentType{
				CategoryPurpose: (*pacs_008_001_08.ExternalCategoryPurpose1Code)(strPtr("CONS")),
			},
			Amount: pacs.FedNowAmount{
				Text: json.Number("1000.00"),
				Ccy:  "USD",
			},
			SenderDI: pacs.FedNowDepositoryInstitution{
				SenderABANumber: "121182904",
			},
			ReceiverDI: pacs.FedNowDepositoryInstitution{
				ReceiverABANumber: "084106768",
			},
			Originator: pacs.FedNowParty{
				Personal: pacs.FedNowPersonal{
					Name: strPtr140("JANE SMITH"),
					Address: pacs.FedNowPstlAdr{
						StreetName:         strPtr70("Dream Road"),
						TownName:           strPtr("Lisle"),
						CountrySubdivision: strPtr("IL"),
						PostalCode:         strPtr16("60532"),
						Country:            countryCodePtr("US"),
					},
				},
			},
			Beneficiary: pacs.FedNowParty{
				Personal: pacs.FedNowPersonal{
					Name: strPtr140("JOHN DOE"),
					Address: pacs.FedNowPstlAdr{
						StreetName:         strPtr70("Dream Road"),
						TownName:           strPtr("Lisle"),
						CountrySubdivision: strPtr("IL"),
						PostalCode:         strPtr16("60532"),
						Country:            countryCodePtr("US"),
					},
				},
			},
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Use a message ID known to be in the full XSD.
		_, err := fednow.Generate(xsdPath, "pacs.008.001.08", cfg, message)
		if err != nil {
			b.Fatal(err)
		}
	}
}
