package admi

import (
	"encoding/xml"

	admi "github.com/mbanq/iso20022-go/ISO20022/admi_002_001_01"
	head "github.com/mbanq/iso20022-go/ISO20022/head_001_001_02"
	"github.com/mbanq/iso20022-go/pkg/common"
	"github.com/mbanq/iso20022-go/pkg/fednow/config"
)

// ParseAdmi002Struct parses an incoming admi.002 MessageReject.
func ParseAdmi002Struct(admiDoc *admi.Document, appHdr head.BusinessApplicationHeaderV02) (*FedNowMessageADM, error) {
	fedMsg := &FedNowMessageADM{
		FedNowMsg: FedNowADM{
			CreationDateTime: common.ISODateTime(appHdr.CreDt),
			Identifier: FedNowIdentifier{
				BusinessMessageID: appHdr.BizMsgIdr,
				MessageType:       appHdr.MsgDefIdr,
				MessageID:         appHdr.BizMsgIdr,
			},
			Reference: admiDoc.Admi00200101.RltdRef.Ref,
			Reason: RejectionReason{
				RejectionReason:   admiDoc.Admi00200101.Rsn.RjctgPtyRsn,
				RejectionDateTime: admiDoc.Admi00200101.Rsn.RjctnDtTm,
			},
		},
	}
	return fedMsg, nil
}

// BuildAdmi002Struct creates an outgoing admi.002 MessageReject.
func BuildAdmi002Struct(message FedNowMessageADM, msgConfig *config.Config) (*admi.Document, error) {
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
