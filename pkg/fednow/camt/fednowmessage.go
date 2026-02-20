package camt

import (
	camt_026_001_07 "github.com/mbanq/iso20022-go/ISO20022/camt_026_001_07"
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

// FlowType constants identify the FedNow business flow context.
// A single ISO message type (e.g. camt.029.001.09) can appear in
// multiple FedNow flows, each requiring a different XML wrapper element.
// The caller sets FlowType on the message so the generator picks the
// correct wrapper. When camt.026/camt.028 support is added later, the
// information-request flow can use FlowTypeInformationRequest.
const (
	// FlowTypeReturnRequest indicates a response to a return/cancellation
	// request (camt.056). Wrapper: FedNowReturnRequestResponse.
	FlowTypeReturnRequest = "return_request"

	// FlowTypeInformationRequest indicates a response to an information
	// request (camt.026/camt.028). Wrapper: FedNowInformationRequestResponse.
	FlowTypeInformationRequest = "information_request"

	// FlowTypeRFPCancellation indicates an RFP cancellation request response.
	// BAH market practice ID: frb.fednow.rcr.01.
	FlowTypeRFPCancellation = "rfp_cancellation"
)

// wrapperByFlowType maps a FlowType to the FedNow envelope wrapper element name.
var wrapperByFlowType = map[string]string{
	FlowTypeReturnRequest:      "FedNowReturnRequestResponse",
	FlowTypeInformationRequest: "FedNowInformationRequestResponse",
}

// marketPracticeIDByFlowType maps camt.029 FlowType to BAH MktPrctc Id per FedNow docs:
// return request response -> frb.fednow.rrr.01, information request response -> frb.fednow.irr.01,
// RFP cancellation request response -> frb.fednow.rcr.01.
var marketPracticeIDByFlowType = map[string]string{
	FlowTypeReturnRequest:      "frb.fednow.rrr.01",
	FlowTypeInformationRequest: "frb.fednow.irr.01",
	FlowTypeRFPCancellation:    "frb.fednow.rcr.01",
}

// ---------------------------------------------------------------------------
// camt.026 – Information Request (Unable To Apply)
// ---------------------------------------------------------------------------

// FedNowMessageInfoReq represents a FedNow camt.026 information request message.
type FedNowMessageInfoReq struct {
	FedNowMsg FedNowInfoReq `json:"fedNowMessage"`
}

func (f FedNowMessageInfoReq) IsFedNowMessage() {}

// PreferredWrapper returns the FedNow envelope wrapper element name for camt.026.
func (f FedNowMessageInfoReq) PreferredWrapper() string {
	return "FedNowInformationRequest"
}

type FedNowInfoReq struct {
	CreationDateTime common.ISODateTime              `json:"creationDateTime"`
	Identifier       FedNowIdentifierInfoReq         `json:"identifier"`
	Case             FedNowCaseInfoReq               `json:"case"`
	Underlying       FedNowUnderlyingInfoReq         `json:"underlying"`
	Justification    FedNowJustification             `json:"justification"`
	SenderDI         FedNowDepositoryInstitutionInfo  `json:"senderDepositoryInstitution"`
	ReceiverDI       FedNowDepositoryInstitutionInfo  `json:"receiverDepositoryInstitution"`
}

type FedNowIdentifierInfoReq struct {
	BusinessMessageID camt_026_001_07.Max35Text `json:"businessMessageId"`
	MessageID         camt_026_001_07.Max35Text `json:"messageId"`
	MessageType       camt_026_001_07.Max35Text `json:"messageType,omitempty"`
	CreationDateTime  common.ISODateTime        `json:"creationDateTime,omitempty"`
}

type FedNowDepositoryInstitutionInfo struct {
	SenderABANumber   camt_026_001_07.Max35Text `json:"senderABANumber,omitempty"`
	ReceiverABANumber camt_026_001_07.Max35Text `json:"receiverABANumber,omitempty"`
}

type FedNowCaseInfoReq struct {
	CaseID    camt_026_001_07.Max35Text       `json:"caseId"`
	CreatorDI FedNowDepositoryInstitutionInfo  `json:"creatorDepositoryInstitution"`
}

type FedNowUnderlyingInfoReq struct {
	Initiation *FedNowUnderlyingInitiation `json:"initiation,omitempty"`
	Interbank  *FedNowUnderlyingInterbank  `json:"interbank,omitempty"`
}

type FedNowUnderlyingInitiation struct {
	OriginalGroupInfo        *FedNowOriginalGroupInfoCamt026                    `json:"originalGroupInformation,omitempty"`
	OriginalPaymentInfoID    *camt_026_001_07.Max35Text                         `json:"originalPaymentInformationId,omitempty"`
	OriginalInstructionID    *camt_026_001_07.Max35Text                         `json:"originalInstructionId,omitempty"`
	OriginalEndToEndID       *camt_026_001_07.Max35Text                         `json:"originalEndToEndId,omitempty"`
	OriginalUETR             *camt_026_001_07.UUIDv4Identifier                  `json:"originalUetr,omitempty"`
	OriginalInstructedAmount *camt_026_001_07.ActiveOrHistoricCurrencyAndAmount `json:"originalInstructedAmount"`
	RequestedExecutionDate   *FedNowDateChoice                                  `json:"requestedExecutionDate,omitempty"`
}

type FedNowUnderlyingInterbank struct {
	OriginalGroupInfo                 *FedNowOriginalGroupInfoCamt026                    `json:"originalGroupInformation,omitempty"`
	OriginalInstructionID             *camt_026_001_07.Max35Text                         `json:"originalInstructionId,omitempty"`
	OriginalEndToEndID                *camt_026_001_07.Max35Text                         `json:"originalEndToEndId,omitempty"`
	OriginalTransactionID             *camt_026_001_07.Max35Text                         `json:"originalTransactionId,omitempty"`
	OriginalUETR                      *camt_026_001_07.UUIDv4Identifier                  `json:"originalUetr,omitempty"`
	OriginalInterbankSettlementAmount *camt_026_001_07.ActiveOrHistoricCurrencyAndAmount `json:"originalInterbankSettlementAmount"`
	OriginalInterbankSettlementDate   *common.ISODate                                    `json:"originalInterbankSettlementDate"`
}

type FedNowOriginalGroupInfoCamt026 struct {
	MessageID        camt_026_001_07.Max35Text `json:"originalMessageId"`
	MessageType      camt_026_001_07.Max35Text `json:"originalMessageType"`
	CreationDateTime common.ISODateTime        `json:"originalCreationDateTime,omitempty"`
}

type FedNowDateChoice struct {
	Date     *common.ISODate     `json:"date,omitempty"`
	DateTime *common.ISODateTime `json:"dateTime,omitempty"`
}

type FedNowJustification struct {
	MissingOrIncorrectInfo       *FedNowMissingOrIncorrectInfo       `json:"missingOrIncorrectInformation,omitempty"`
	PossibleDuplicateInstruction *camt_026_001_07.TrueFalseIndicator `json:"possibleDuplicateInstruction,omitempty"`
}

type FedNowMissingOrIncorrectInfo struct {
	AMLRequest    *camt_026_001_07.AMLIndicator `json:"amlRequest,omitempty"`
	MissingInfo   []FedNowMissingInfo           `json:"missingInformation,omitempty"`
	IncorrectInfo []FedNowIncorrectInfo         `json:"incorrectInformation,omitempty"`
}

type FedNowMissingInfo struct {
	Code           camt_026_001_07.UnableToApplyMissingInformation3Code `json:"code"`
	AdditionalInfo *camt_026_001_07.Max140Text                         `json:"additionalInformation,omitempty"`
}

type FedNowIncorrectInfo struct {
	Code           camt_026_001_07.UnableToApplyIncorrectInformation4Code `json:"code"`
	AdditionalInfo *camt_026_001_07.Max140Text                           `json:"additionalInformation,omitempty"`
}

// ---------------------------------------------------------------------------
// camt.029 – Resolution of Investigation
// ---------------------------------------------------------------------------

// FedNowMessageCxlRsp represents a FedNow camt.029 cancellation response message.
// It implements fednow.FedNowMessage via IsFedNowMessage().
type FedNowMessageCxlRsp struct {
	FedNowMsg FedNowCxlRsp `json:"fedNowMessage"`
}

func (f FedNowMessageCxlRsp) IsFedNowMessage() {}

// PreferredWrapper returns the FedNow envelope wrapper element name based on
// FlowType. Defaults to FedNowReturnRequestResponse for backward compatibility
// when FlowType is not set.
func (f FedNowMessageCxlRsp) PreferredWrapper() string {
	if wrapper, ok := wrapperByFlowType[f.FedNowMsg.FlowType]; ok {
		return wrapper
	}
	// Default: return-request flow (the only currently supported flow).
	return "FedNowReturnRequestResponse"
}

// MarketPracticeID returns the BAH MktPrctc Id for this camt.029 message based on
// FlowType. Defaults to frb.fednow.rrr.01 (return request response) when FlowType is not set.
func (f FedNowMessageCxlRsp) MarketPracticeID() string {
	if id, ok := marketPracticeIDByFlowType[f.FedNowMsg.FlowType]; ok {
		return id
	}
	return "frb.fednow.rrr.01"
}

// FedNowCxlRsp is the custom JSON payload used by this library for camt.029.
type FedNowCxlRsp struct {
	// FlowType identifies the business flow context (e.g. FlowTypeReturnRequest
	// or FlowTypeInformationRequest). This determines the FedNow envelope wrapper
	// element used during XML generation.
	FlowType            string                       `json:"flowType,omitempty"`
	CreationDateTime    common.ISODateTime           `json:"creationDateTime"`
	Identifier          FedNowIdentifierCxlRsp       `json:"identifier"`
	ResolvedCase        FedNowCase                   `json:"resolvedCase"`
	InvestigationStatus FedNowInvestigationStatus    `json:"investigationStatus"`
	CancellationDetails []FedNowCxlRspDetails        `json:"cancellationDetails,omitempty"`
	SenderDI            FedNowDepositoryInstitution2 `json:"senderDepositoryInstitution"`
	ReceiverDI          FedNowDepositoryInstitution2 `json:"receiverDepositoryInstitution"`

	// IRR-specific fields (Information Request Response flow).
	// CorrectionTransaction conveys interbank payment reference when Status = IPAY.
	CorrectionTransaction *FedNowCorrectionTransaction `json:"correctionTransaction,omitempty"`
	// ResolutionRelatedInformation conveys original payment identifiers when Status = IDUP.
	ResolutionRelatedInformation *FedNowResolutionRelatedInfo `json:"resolutionRelatedInformation,omitempty"`
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

// FedNowCorrectionTransaction represents the Interbank correction transaction
// (CrrctnTx/IntrBk) used in Information Request Response when Status = IPAY.
type FedNowCorrectionTransaction struct {
	GroupHeader               FedNowCorrectionGroupHeader                        `json:"groupHeader"`
	InstructionID             *camt_029_001_09.Max35Text                         `json:"instructionId,omitempty"`
	EndToEndID                *camt_029_001_09.Max35Text                         `json:"endToEndId,omitempty"`
	TransactionID             *camt_029_001_09.Max35Text                         `json:"transactionId,omitempty"`
	UETR                      *camt_029_001_09.UUIDv4Identifier                  `json:"uetr,omitempty"`
	InterbankSettlementAmount camt_029_001_09.ActiveOrHistoricCurrencyAndAmount  `json:"interbankSettlementAmount"`
	InterbankSettlementDate   common.ISODate                                     `json:"interbankSettlementDate"`
}

// FedNowCorrectionGroupHeader maps to CorrectiveGroupInformation1 (GrpHdr
// within CorrectiveInterbankTransaction2).
type FedNowCorrectionGroupHeader struct {
	MessageID        camt_029_001_09.Max35Text `json:"messageId"`
	MessageNameID    camt_029_001_09.Max35Text `json:"messageNameId"`
	CreationDateTime *common.ISODateTime       `json:"creationDateTime,omitempty"`
}
