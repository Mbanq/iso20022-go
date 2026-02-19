package pacs

import (
	"encoding/json"
	"encoding/xml"

	head "github.com/mbanq/iso20022-go/ISO20022/head_001_001_02"
	"github.com/mbanq/iso20022-go/ISO20022/pacs_008_001_08"
	"github.com/mbanq/iso20022-go/ISO20022/pa
	"github.com/mbanq/iso20022-go/pkg/common"
	"github.com/mbanq/iso20022-go/pkg/fednow/config"
)

// FedNowMessageStsReq represents a FedNow pacs.028 payment status request message.
type FedNowMessageStsReq struct {
	FedNowMsg FedNowStsReq `json:"fedNowMessage"`
}

func (f FedNowMessageStsReq) IsFedNowMessage() {}

// FedNowStsReq is the custom JSON payload for pacs.028 status requests.
type FedNowStsReq struct {
	CreationDateTime   common.ISODateTime          `json:"creationDateTime"`
	Identifier         FedNowIdentifier            `json:"identifier"`
	OriginalIdentifier FedNowIdentifier            `json:"originalIdentifier"`
	SenderDI           FedNowDepositoryInstitution `json:"senderDepositoryInstitution"`
	ReceiverDI         FedNowDepositoryInstitution `json:"receiverDepositoryInstitution"`
}

// ParsePacs028 parses an incoming pacs.028 payment status request.
func ParsePacs028(appHdr head.BusinessApplicationHeaderV02, document pacs_028_001_03.Document) (*FedNowMessageStsReq, error) {

	fitofipmtstsreq := document.FIToFIPmtStsReq

	// Extract transaction info (first transaction in the array)
	var txInfo pacs_028_001_03.PaymentTransaction113
	if len(fitofipmtstsreq.TxInf) > 0 {
		txInfo = fitofipmtstsreq.TxInf[0]
	}

	// Extract original message identifiers
	var orgnlMsgId pacs_008_001_08.Max35Text
	var orgnlMsgNmId pacs_008_001_08.Max35Text
	var orgnlCreDtTm common.ISODateTime
	if txInfo.OrgnlGrpInf != nil {
		orgnlMsgId = pacs_008_001_08.Max35Text(txInfo.OrgnlGrpInf.OrgnlMsgId)
		orgnlMsgNmId = pacs_008_001_08.Max35Text(txInfo.OrgnlGrpInf.OrgnlMsgNmId)
		if txInfo.OrgnlGrpInf.OrgnlCreDtTm != nil {
			orgnlCreDtTm = common.ISODateTime(*txInfo.OrgnlGrpInf.OrgnlCreDtTm)
		}
	}

	// Extract original instruction ID
	var orgnlInstrId *pacs_008_001_08.Max35Text
	if txInfo.OrgnlInstrId != nil {
		val := pacs_008_001_08.Max35Text(*txInfo.OrgnlInstrId)
		orgnlInstrId = &val
	}

	// Extract original end-to-end ID
	var orgnlEndToEndId pacs_008_001_08.Max35Text
	if txInfo.OrgnlEndToEndId != nil {
		orgnlEndToEndId = pacs_008_001_08.Max35Text(*txInfo.OrgnlEndToEndId)
	}

	// Extract original transaction ID
	var orgnlTxId *pacs_008_001_08.Max35Text
	if txInfo.OrgnlTxId != nil {
		val := pacs_008_001_08.Max35Text(*txInfo.OrgnlTxId)
		orgnlTxId = &val
	}

	// Extract original UETR
	var orgnlUETR *pacs_008_001_08.UUIDv4Identifier
	if txInfo.OrgnlUETR != nil {
		val := pacs_008_001_08.UUIDv4Identifier(*txInfo.OrgnlUETR)
		orgnlUETR = &val
	}

	// Extract routing numbers from AppHdr
	senderABANumber := extractClrSysMemberID(appHdr.Fr)
	receiverABANumber := extractClrSysMemberID(appHdr.To)

	fednowMsg := FedNowMessageStsReq{
		FedNowMsg: FedNowStsReq{
			CreationDateTime: common.ISODateTime(fitofipmtstsreq.GrpHdr.CreDtTm),
			Identifier: FedNowIdentifier{
				BusinessMessageID: pacs_008_001_08.Max35Text(appHdr.BizMsgIdr),
				MessageID:         pacs_008_001_08.Max35Text(fitofipmtstsreq.GrpHdr.MsgId),
				CreationDateTime:  common.ISODateTime(appHdr.CreDt),
			},
			OriginalIdentifier: FedNowIdentifier{
				MessageID:        orgnlMsgId,
				MessageType:      orgnlMsgNmId,
				InstructionID:    orgnlInstrId,
				EndToEndID:       orgnlEndToEndId,
				TransactionID:    orgnlTxId,
				UETR:             orgnlUETR,
				CreationDateTime: orgnlCreDtTm,
			},
			SenderDI: FedNowDepositoryInstitution{
				SenderABANumber: senderABANumber,
			},
			ReceiverDI: FedNowDepositoryInstitution{
				ReceiverABANumber: receiverABANumber,
			},
		},
	}

	return &fednowMsg, nil
}

// BuildPacs028Struct creates a pacs.028 Document from a FedNowMessageStsReq.
// This is for future outbound status requests (not currently needed).
func BuildPacs028Struct(message FedNowMessageStsReq, msgConfig *config.Config) (*pacs_028_001_03.Document, error) {

	fedMsg := message.FedNowMsg

	clearingSystemId := pacs_028_001_03.ExternalClearingSystemIdentification1Code(msgConfig.ClearingSystemId)

	// Build OriginalGroupInformation if we have original message details
	var orgnlGrpInf *pacs_028_001_03.OriginalGroupInformation29
	if fedMsg.OriginalIdentifier.MessageID != "" {
		orgnlCreDtTm := fedMsg.OriginalIdentifier.CreationDateTime
		var creationTimePtr *common.ISODateTime
		if orgnlCreDtTm != (common.ISODateTime{}) {
			creationTimePtr = &orgnlCreDtTm
		}

		orgnlGrpInf = &pacs_028_001_03.OriginalGroupInformation29{
			OrgnlMsgId:   pacs_028_001_03.Max35Text(fedMsg.OriginalIdentifier.MessageID),
			OrgnlMsgNmId: pacs_028_001_03.Max35Text(fedMsg.OriginalIdentifier.MessageType),
			OrgnlCreDtTm: creationTimePtr,
		}
	}

	// Build optional fields
	var orgnlInstrId *pacs_028_001_03.Max35Text
	if fedMsg.OriginalIdentifier.InstructionID != nil {
		val := pacs_028_001_03.Max35Text(*fedMsg.OriginalIdentifier.InstructionID)
		orgnlInstrId = &val
	}

	var orgnlEndToEndId *pacs_028_001_03.Max35Text
	if fedMsg.OriginalIdentifier.EndToEndID != "" {
		val := pacs_028_001_03.Max35Text(fedMsg.OriginalIdentifier.EndToEndID)
		orgnlEndToEndId = &val
	}

	var orgnlTxId *pacs_028_001_03.Max35Text
	if fedMsg.OriginalIdentifier.TransactionID != nil {
		val := pacs_028_001_03.Max35Text(*fedMsg.OriginalIdentifier.TransactionID)
		orgnlTxId = &val
	}

	var orgnlUETR *pacs_028_001_03.UUIDv4Identifier
	if fedMsg.OriginalIdentifier.UETR != nil {
		val := pacs_028_001_03.UUIDv4Identifier(*fedMsg.OriginalIdentifier.UETR)
		orgnlUETR = &val
	}

	pacsDoc := &pacs_028_001_03.Document{
		XMLName: xml.Name{Space: "urn:iso:std:iso:20022:tech:xsd:pacs.028.001.03", Local: "Document"},
		FIToFIPmtStsReq: pacs_028_001_03.FIToFIPaymentStatusRequestV03{
			GrpHdr: pacs_028_001_03.GroupHeader91{
				MsgId:   pacs_028_001_03.Max35Text(fedMsg.Identifier.MessageID),
				CreDtTm: (common.ISODateTime)(fedMsg.CreationDateTime),
				InstgAgt: &pacs_028_001_03.BranchAndFinancialInstitutionIdentification6{
					FinInstnId: pacs_028_001_03.FinancialInstitutionIdentification18{
						ClrSysMmbId: &pacs_028_001_03.ClearingSystemMemberIdentification2{
							MmbId: pacs_028_001_03.Max35Text(fedMsg.SenderDI.SenderABANumber),
							ClrSysId: &pacs_028_001_03.ClearingSystemIdentification2Choice{
								Cd: &clearingSystemId,
							},
						},
					},
				},
				InstdAgt: &pacs_028_001_03.BranchAndFinancialInstitutionIdentification6{
					FinInstnId: pacs_028_001_03.FinancialInstitutionIdentification18{
						ClrSysMmbId: &pacs_028_001_03.ClearingSystemMemberIdentification2{
							MmbId: pacs_028_001_03.Max35Text(fedMsg.ReceiverDI.ReceiverABANumber),
							ClrSysId: &pacs_028_001_03.ClearingSystemIdentification2Choice{
								Cd: &clearingSystemId,
							},
						},
					},
				},
			},
			TxInf: []pacs_028_001_03.PaymentTransaction113{
				{
					OrgnlGrpInf:     orgnlGrpInf,
					OrgnlInstrId:    orgnlInstrId,
					OrgnlEndToEndId: orgnlEndToEndId,
					OrgnlTxId:       orgnlTxId,
					OrgnlUETR:       orgnlUETR,
				},
			},
		},
	}

	return pacsDoc, nil
}

// BuildPacs028 creates a pacs.028 Document from JSON payload.
func BuildPacs028(payload []byte, config *config.Config) (*pacs_028_001_03.Document, error) {
	var message FedNowMessageStsReq
	if err := json.Unmarshal(payload, &message); err != nil {
		return nil, err
	}

	return BuildPacs028Struct(message, config)
}
