package camt

import (
	camt_029_001_09 "github.com/mbanq/iso20022-go/ISO20022/camt_029_001_09"
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

// FedNowMessageCxlRsp represents a FedNow camt.029 cancellation response message.
// It implements fednow.FedNowMessage via IsFedNowMessage().
type FedNowMessageCxlRsp struct {
	FedNowMsg FedNowCxlRsp `json:"fedNowMessage"`
}

func (f FedNowMessageCxlRsp) IsFedNowMessage() {}

// FedNowCxlRsp is the custom JSON payload used by this library for camt.029.
type FedNowCxlRsp struct {
	CreationDateTime    common.ISODateTime           `json:"creationDateTime"`
	Identifier          FedNowIdentifierCxlRsp       `json:"identifier"`
	ResolvedCase        FedNowCase                   `json:"resolvedCase"`
	InvestigationStatus FedNowInvestigationStatus    `json:"investigationStatus"`
	CancellationDetails []FedNowCxlRspDetails        `json:"cancellationDetails,omitempty"`
	SenderDI            FedNowDepositoryInstitution2 `json:"senderDepositoryInstitution"`
	ReceiverDI          FedNowDepositoryInstitution2 `json:"receiverDepositoryInstitution"`
}

type FedNowIdentifierCxlRsp struct {
	BusinessMessageID camt_029_001_09.Max35Text `json:"businessMessageId"`
	MessageID         camt_029_001_09.Max35Text `json:"messageId"`
	MessageType       camt_029_001_09.Max35Text `json:"messageType,omitempty"`
	CreationDateTime  common.ISODateTime        `json:"creationDateTime,omitempty"`
}

type FedNowDepositoryInstitution2 struct {
	SenderABANumber   camt_029_001_09.Max35Text   `json:"senderABANumber,omitempty"`
	ReceiverABANumber camt_029_001_09.Max35Text   `json:"receiverABANumber,omitempty"`
	Name              *camt_029_001_09.Max140Text `json:"senderShortName,omitempty"`
}

type FedNowCase struct {
	CaseID    camt_029_001_09.Max35Text    `json:"caseId"`
	CreatorDI FedNowDepositoryInstitution2 `json:"creatorDepositoryInstitution"`
}

type FedNowInvestigationStatus struct {
	Confirmation                    *camt_029_001_09.ExternalInvestigationExecutionConfirmation1Code `json:"confirmation,omitempty"`
	RejectedModification            []camt_029_001_09.ModificationStatusReason1Choice                `json:"rejectedModification,omitempty"`
	DuplicateOf                     *FedNowCase                                                      `json:"duplicateOf,omitempty"`
	AssignmentCancellationConfirmed *camt_029_001_09.YesNoIndicator                                  `json:"assignmentCancellationConfirmed,omitempty"`
}

type FedNowOriginalGroupInfo struct {
	MessageID        camt_029_001_09.Max35Text `json:"originalMessageId"`
	MessageType      camt_029_001_09.Max35Text `json:"originalMessageType"`
	CreationDateTime common.ISODateTime        `json:"originalCreationDateTime,omitempty"`
}

type FedNowResolutionRelatedInfo struct {
	EndToEndID                *camt_029_001_09.Max35Text                         `json:"endToEndId,omitempty"`
	TransactionID             *camt_029_001_09.Max35Text                         `json:"transactionId,omitempty"`
	UETR                      *camt_029_001_09.UUIDv4Identifier                  `json:"uetr,omitempty"`
	InterbankSettlementAmount *camt_029_001_09.ActiveOrHistoricCurrencyAndAmount `json:"interbankSettlementAmount,omitempty"`
	InterbankSettlementDate   *common.ISODate                                    `json:"interbankSettlementDate,omitempty"`
}

type FedNowCxlRspDetails struct {
	OriginalGroupInfo     *FedNowOriginalGroupInfo          `json:"originalGroupInformation,omitempty"`
	OriginalInstructionID *camt_029_001_09.Max35Text        `json:"originalInstructionId,omitempty"`
	OriginalEndToEndID    *camt_029_001_09.Max35Text        `json:"originalEndToEndId,omitempty"`
	OriginalUETR          *camt_029_001_09.UUIDv4Identifier `json:"originalUetr,omitempty"`
	ResolutionRelatedInfo *FedNowResolutionRelatedInfo      `json:"resolutionRelatedInformation,omitempty"`
}
