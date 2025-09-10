package admi

import (
	"github.com/Mbanq/iso20022-go/ISO20022/admi_002_001_01"
	head "github.com/Mbanq/iso20022-go/ISO20022/head_001_001_02"
	"github.com/Mbanq/iso20022-go/pkg/common"
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

type FedNowIdentifier struct {
	BusinessMessageID head.Max35Text `json:"businessMessageId"`
	MessageType       head.Max35Text `json:"messageType"`
	MessageID         head.Max35Text `json:"messageId"`
}
