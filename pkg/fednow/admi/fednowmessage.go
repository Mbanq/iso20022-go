package admi

import (
	admi_002_001_01 "github.com/mbanq/iso20022-go/ISO20022/admi_002_001_01"
	admi_007_001_01 "github.com/mbanq/iso20022-go/ISO20022/admi_007_001_01"
	head "github.com/mbanq/iso20022-go/ISO20022/head_001_001_02"
	"github.com/mbanq/iso20022-go/pkg/common"
)

type FedNowMessageADM struct {
	FedNowMsg FedNowADM `json:"fedNowMessage"`
}

func (f FedNowMessageADM) IsFedNowMessage() {}

type FedNowADM struct {
	CreationDateTime common.ISODateTime        `json:"creationDateTime"`
	Identifier       FedNowIdentifier          `json:"identifier"`
	Reference        admi_002_001_01.Max35Text `json:"reference"`
	Reason           RejectionReason           `json:"reason"`
}

type RejectionReason struct {
	RejectionReason   admi_002_001_01.Max35Text `json:"rejectionReason"`
	RejectionDateTime *common.ISODateTime       `json:"rejectionDateTime"`
}

type FedNowMessageRctAck struct {
	FedNowMsg FedNowReceiptAcknowledgement `json:"fedNowMessage"`
}

func (f FedNowMessageRctAck) IsFedNowMessage() {}

type FedNowReceiptAcknowledgement struct {
	CreationDateTime common.ISODateTime             `json:"creationDateTime"`
	Identifier       FedNowIdentifier               `json:"identifier"`
	QueryName        *admi_007_001_01.Max35Text     `json:"queryName,omitempty"`
	Reports          []ReceiptAcknowledgementReport `json:"reports"`
}

type ReceiptAcknowledgementReport struct {
	RelatedReference ReceiptAcknowledgementReference `json:"relatedReference"`
	RequestHandling  ReceiptAcknowledgementHandling  `json:"requestHandling"`
}

type ReceiptAcknowledgementReference struct {
	Reference       admi_007_001_01.Max35Text               `json:"reference"`
	MessageName     *admi_007_001_01.Max35Text              `json:"messageName,omitempty"`
	ReferenceIssuer *admi_007_001_01.PartyIdentification136 `json:"referenceIssuer,omitempty"`
}

type ReceiptAcknowledgementHandling struct {
	StatusCode     admi_007_001_01.Max4AlphaNumericText `json:"statusCode"`
	StatusDateTime *common.ISODateTime                  `json:"statusDateTime,omitempty"`
	Description    *admi_007_001_01.Max140Text          `json:"description,omitempty"`
}

type FedNowIdentifier struct {
	BusinessMessageID head.Max35Text `json:"businessMessageId"`
	MessageType       head.Max35Text `json:"messageType"`
	MessageID         head.Max35Text `json:"messageId"`
}
