package admi

import (
	admi_002_001_01 "github.com/mbanq/iso20022-go/ISO20022/admi_002_001_01"
	admi_007_001_01 "github.com/mbanq/iso20022-go/ISO20022/admi_007_001_01"
	head "github.com/mbanq/iso20022-go/ISO20022/head_001_001_02"
	"github.com/mbanq/iso20022-go/pkg/common"
)

// FedNowMessageADM represents an admi.002 MessageReject.
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

// FedNowMessageParticipantBroadcast represents an outgoing admi.004 ParticipantBroadcast.
// Used for FPON (sign on), FPOF (sign off), PING, FPCR (reconnect), FPCD (disconnect).
type FedNowMessageParticipantBroadcast struct {
	FedNowMsg FedNowParticipantBroadcast `json:"fedNowMessage"`
}

func (f FedNowMessageParticipantBroadcast) IsFedNowMessage() {}

type FedNowParticipantBroadcast struct {
	MessageID            string              `json:"messageId"`
	EventCode            string              `json:"eventCode"`                      // FPON, FPOF, PING, FPCR, FPCD
	ConnectionIdentifier string              `json:"connectionIdentifier,omitempty"` // RTN or Connection Point ID (7-9 chars)
	EventDescription     *string             `json:"eventDescription,omitempty"`
	EventTime            *common.ISODateTime `json:"eventTime,omitempty"`
}

// FedNowMessageFedNowBroadcast represents an incoming admi.004 FedNowBroadcast.
// Sent by FedNow to all participants to notify of system or participant events.
type FedNowMessageFedNowBroadcast struct {
	FedNowMsg FedNowFedNowBroadcast `json:"fedNowMessage"`
}

func (f FedNowMessageFedNowBroadcast) IsFedNowMessage() {}

type FedNowFedNowBroadcast struct {
	CreationDateTime common.ISODateTime  `json:"creationDateTime"`
	MessageID        string              `json:"messageId"`
	EventCode        string              `json:"eventCode"`             // FPON, FPOF, ROLL, EXTN, FNKY, etc.
	EventParams      []string            `json:"eventParams,omitempty"` // Multiple parameters allowed
	EventDescription *string             `json:"eventDescription,omitempty"`
	EventTime        *common.ISODateTime `json:"eventTime,omitempty"`
}

// FedNowMessageSystemResponse represents an incoming admi.011 FedNowSystemResponse.
// FedNow's response to our admi.004 ParticipantBroadcast (confirmation of FPON/FPOF/PING).
type FedNowMessageSystemResponse struct {
	FedNowMsg FedNowSystemResponse `json:"fedNowMessage"`
}

func (f FedNowMessageSystemResponse) IsFedNowMessage() {}

type FedNowSystemResponse struct {
	CreationDateTime  common.ISODateTime  `json:"creationDateTime"`
	MessageID         string              `json:"messageId"`
	OriginalReference string              `json:"originalReference"`    // Our original admi.004 message ID
	EventCode         string              `json:"eventCode,omitempty"`  // FPON, FPOF, PING, etc.
	EventParam        string              `json:"eventParam,omitempty"` // Connection identifier
	EventTime         *common.ISODateTime `json:"eventTime,omitempty"`
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
