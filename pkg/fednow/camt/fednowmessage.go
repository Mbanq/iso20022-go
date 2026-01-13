package camt

import (
	camt_056_001_08 "github.com/mbanq/iso20022-go/ISO20022/camt_056_001_08"
	"github.com/mbanq/iso20022-go/pkg/common"
)

// FedNowMessageCxlReq represents a FedNow camt.056 cancellation request message.
// It implements fednow.FedNowMessage via IsFedNowMessage().
type FedNowMessageCxlReq struct {
	FedNowMsg FedNowCxlReq `json:"fedNowMessage"`
}

func (f FedNowMessageCxlReq) IsFedNowMessage() {}

// FedNowCxlReq is the custom JSON payload used by this library for camt.056.
// It intentionally mirrors the style used in other message folders (pacs/pain).
type FedNowCxlReq struct {
	CreationDateTime   common.ISODateTime                               `json:"creationDateTime"`
	Identifier         FedNowIdentifier                                 `json:"identifier"`
	OriginalIdentifier FedNowIdentifier                                 `json:"originalIdentifier"`
	CancellationReason *camt_056_001_08.ExternalCancellationReason1Code `json:"cancellationReason,omitempty"`
	AdditionalInfo     *camt_056_001_08.Max105Text                      `json:"additionalInformation,omitempty"`
	SenderDI           FedNowDepositoryInstitution                      `json:"senderDepositoryInstitution"`
	ReceiverDI         FedNowDepositoryInstitution                      `json:"receiverDepositoryInstitution"`
}

type FedNowIdentifier struct {
	BusinessMessageID camt_056_001_08.Max35Text         `json:"businessMessageId"`
	MessageID         camt_056_001_08.Max35Text         `json:"messageId"`
	MessageType       camt_056_001_08.Max35Text         `json:"messageType,omitempty"`
	InstructionID     *camt_056_001_08.Max35Text        `json:"instructionId,omitempty"`
	EndToEndID        camt_056_001_08.Max35Text         `json:"endToEndId,omitempty"`
	TransactionID     *camt_056_001_08.Max35Text        `json:"transactionId,omitempty"`
	UETR              *camt_056_001_08.UUIDv4Identifier `json:"uetr,omitempty"`
	CreationDateTime  common.ISODateTime                `json:"creationDateTime,omitempty"`
}

type FedNowDepositoryInstitution struct {
	SenderABANumber   camt_056_001_08.Max35Text   `json:"senderABANumber,omitempty"`
	ReceiverABANumber camt_056_001_08.Max35Text   `json:"receiverABANumber,omitempty"`
	Name              *camt_056_001_08.Max140Text `json:"senderShortName,omitempty"`
}
