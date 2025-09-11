package admi

import (
	"encoding/xml"

	admi "github.com/mbanq/iso20022-go/ISO20022/admi_002_001_01"
	"github.com/mbanq/iso20022-go/pkg/fednow/config"
)

func BuildAdmi004Struct(message FedNowMessageADM, msgConfig *config.Config) (*admi.Document, error) {

	fedMsg := message.FedNowMsg

	admiDoc := &admi.Document{
		XMLName: xml.Name{Space: "urn:iso:std:iso:20022:tech:xsd:admi.002.001.01", Local: "Document"},
		Admi00200101: admi.Admi00200101{
			RltdRef: admi.MessageReference{
				Ref: admi.Max35Text(fedMsg.Reference),
			},
			Rsn: admi.RejectionReason2{
				RjctgPtyRsn: admi.Max35Text(fedMsg.Reason.RejectionReason),
				RjctnDtTm:   fedMsg.Reason.RejectionDateTime,
			},
		},
	}
	return admiDoc, nil
}
