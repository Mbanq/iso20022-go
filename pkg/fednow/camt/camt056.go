package camt

import (
	"encoding/json"
	"encoding/xml"
	"strings"
	"time"

	camt_056_001_08 "github.com/mbanq/iso20022-go/ISO20022/camt_056_001_08"
	head "github.com/mbanq/iso20022-go/ISO20022/head_001_001_02"
	"github.com/mbanq/iso20022-go/pkg/common"
	"github.com/mbanq/iso20022-go/pkg/fednow/config"
)

func BuildCamt056Struct(message FedNowMessageCxlReq, msgConfig *config.Config) (*camt_056_001_08.Document, error) {
	fedMsg := message.FedNowMsg

	clearingSystemId := camt_056_001_08.ExternalClearingSystemIdentification1Code(msgConfig.ClearingSystemId)

	// OrgnlCreDtTm is optional.
	var orgnlCreationTime *common.ISODateTime
	if !time.Time(fedMsg.OriginalIdentifier.CreationDateTime).IsZero() {
		value := fedMsg.OriginalIdentifier.CreationDateTime
		orgnlCreationTime = &value
	}

	assgnr := camt_056_001_08.Party40Choice{
		Agt: &camt_056_001_08.BranchAndFinancialInstitutionIdentification6{
			FinInstnId: camt_056_001_08.FinancialInstitutionIdentification18{
				ClrSysMmbId: &camt_056_001_08.ClearingSystemMemberIdentification2{
					MmbId: camt_056_001_08.Max35Text(fedMsg.SenderDI.SenderABANumber),
					ClrSysId: &camt_056_001_08.ClearingSystemIdentification2Choice{
						Cd: &clearingSystemId,
					},
				},
			},
		},
	}
	assgne := camt_056_001_08.Party40Choice{
		Agt: &camt_056_001_08.BranchAndFinancialInstitutionIdentification6{
			FinInstnId: camt_056_001_08.FinancialInstitutionIdentification18{
				ClrSysMmbId: &camt_056_001_08.ClearingSystemMemberIdentification2{
					MmbId: camt_056_001_08.Max35Text(fedMsg.ReceiverDI.ReceiverABANumber),
					ClrSysId: &camt_056_001_08.ClearingSystemIdentification2Choice{
						Cd: &clearingSystemId,
					},
				},
			},
		},
	}

	// Cancellation reason (optional) mapped onto both group and transaction levels.
	var cxlRsnInf []camt_056_001_08.PaymentCancellationReason5
	if fedMsg.CancellationReason != nil || (fedMsg.AdditionalInfo != nil && strings.TrimSpace(string(*fedMsg.AdditionalInfo)) != "") {
		var rsnChoice *camt_056_001_08.CancellationReason33Choice
		if fedMsg.CancellationReason != nil {
			rsnChoice = &camt_056_001_08.CancellationReason33Choice{
				Cd: fedMsg.CancellationReason,
			}
		}
		var addtl []camt_056_001_08.Max105Text
		if fedMsg.AdditionalInfo != nil && strings.TrimSpace(string(*fedMsg.AdditionalInfo)) != "" {
			addtl = []camt_056_001_08.Max105Text{*fedMsg.AdditionalInfo}
		}

		cxlRsnInf = []camt_056_001_08.PaymentCancellationReason5{
			{
				Rsn:      rsnChoice,
				AddtlInf: addtl,
			},
		}
	}

	origGrp := camt_056_001_08.OriginalGroupHeader15{
		OrgnlMsgId:   camt_056_001_08.Max35Text(fedMsg.OriginalIdentifier.MessageID),
		OrgnlMsgNmId: camt_056_001_08.Max35Text(fedMsg.OriginalIdentifier.MessageType),
		OrgnlCreDtTm: orgnlCreationTime,
		CxlRsnInf:    cxlRsnInf,
	}

	txInf := camt_056_001_08.PaymentTransaction106{
		OrgnlInstrId:    fedMsg.OriginalIdentifier.InstructionID,
		OrgnlEndToEndId: (*camt_056_001_08.Max35Text)(&fedMsg.OriginalIdentifier.EndToEndID),
		OrgnlTxId:       fedMsg.OriginalIdentifier.TransactionID,
		OrgnlUETR:       fedMsg.OriginalIdentifier.UETR,
		Assgnr:          assgnr.Agt,
		Assgne:          assgne.Agt,
		CxlRsnInf:       cxlRsnInf,
	}

	doc := &camt_056_001_08.Document{
		XMLName: xml.Name{Space: "urn:iso:std:iso:20022:tech:xsd:camt.056.001.08", Local: "Document"},
		FIToFIPmtCxlReq: camt_056_001_08.FIToFIPaymentCancellationRequestV08{
			Assgnmt: camt_056_001_08.CaseAssignment5{
				Id:      camt_056_001_08.Max35Text(fedMsg.Identifier.MessageID),
				Assgnr:  assgnr,
				Assgne:  assgne,
				CreDtTm: fedMsg.CreationDateTime,
			},
			Undrlyg: []camt_056_001_08.UnderlyingTransaction23{
				{
					OrgnlGrpInfAndCxl: &origGrp,
					TxInf:             []camt_056_001_08.PaymentTransaction106{txInf},
				},
			},
		},
	}

	return doc, nil
}

func BuildCamt056(payload []byte, cfg *config.Config) (*camt_056_001_08.Document, error) {
	var message FedNowMessageCxlReq
	if err := json.Unmarshal(payload, &message); err != nil {
		return nil, err
	}
	return BuildCamt056Struct(message, cfg)
}

func ParseCamt056(appHdr head.BusinessApplicationHeaderV02, document camt_056_001_08.Document) (*FedNowMessageCxlReq, error) {
	req := document.FIToFIPmtCxlReq

	// Extract original identifiers (best-effort, first underlying + first tx).
	var (
		origMsgId      camt_056_001_08.Max35Text
		origMsgNmId    camt_056_001_08.Max35Text
		origCreDtTm    common.ISODateTime
		origInstrId    *camt_056_001_08.Max35Text
		origEndToEndId camt_056_001_08.Max35Text
		origTxId       *camt_056_001_08.Max35Text
		origUETR       *camt_056_001_08.UUIDv4Identifier
		cxlReason      *camt_056_001_08.ExternalCancellationReason1Code
		addtlInfo      *camt_056_001_08.Max105Text
	)

	if len(req.Undrlyg) > 0 && req.Undrlyg[0].OrgnlGrpInfAndCxl != nil {
		grp := req.Undrlyg[0].OrgnlGrpInfAndCxl
		origMsgId = grp.OrgnlMsgId
		origMsgNmId = grp.OrgnlMsgNmId
		if grp.OrgnlCreDtTm != nil {
			origCreDtTm = common.ISODateTime(*grp.OrgnlCreDtTm)
		}
		if len(grp.CxlRsnInf) > 0 && grp.CxlRsnInf[0].Rsn != nil && grp.CxlRsnInf[0].Rsn.Cd != nil {
			cxlReason = grp.CxlRsnInf[0].Rsn.Cd
		}
		if len(grp.CxlRsnInf) > 0 && len(grp.CxlRsnInf[0].AddtlInf) > 0 {
			tmp := camt_056_001_08.Max105Text(grp.CxlRsnInf[0].AddtlInf[0])
			addtlInfo = &tmp
		}
	}

	if len(req.Undrlyg) > 0 && len(req.Undrlyg[0].TxInf) > 0 {
		tx := req.Undrlyg[0].TxInf[0]
		origInstrId = tx.OrgnlInstrId
		if tx.OrgnlEndToEndId != nil {
			origEndToEndId = *tx.OrgnlEndToEndId
		}
		origTxId = tx.OrgnlTxId
		origUETR = tx.OrgnlUETR
		// If group didn't have reason, fall back to tx-level.
		if cxlReason == nil && len(tx.CxlRsnInf) > 0 && tx.CxlRsnInf[0].Rsn != nil && tx.CxlRsnInf[0].Rsn.Cd != nil {
			cxlReason = tx.CxlRsnInf[0].Rsn.Cd
		}
		if addtlInfo == nil && len(tx.CxlRsnInf) > 0 && len(tx.CxlRsnInf[0].AddtlInf) > 0 {
			tmp := camt_056_001_08.Max105Text(tx.CxlRsnInf[0].AddtlInf[0])
			addtlInfo = &tmp
		}
	}

	senderABANumber := extractClrSysMemberID(appHdr.Fr)
	receiverABANumber := extractClrSysMemberID(appHdr.To)

	msg := FedNowMessageCxlReq{
		FedNowMsg: FedNowCxlReq{
			CreationDateTime: common.ISODateTime(req.Assgnmt.CreDtTm),
			Identifier: FedNowIdentifier{
				BusinessMessageID: camt_056_001_08.Max35Text(appHdr.BizMsgIdr),
				MessageID:         camt_056_001_08.Max35Text(req.Assgnmt.Id),
				MessageType:       camt_056_001_08.Max35Text(appHdr.MsgDefIdr),
				CreationDateTime:  common.ISODateTime(appHdr.CreDt),
			},
			OriginalIdentifier: FedNowIdentifier{
				MessageID:        origMsgId,
				MessageType:      origMsgNmId,
				InstructionID:    origInstrId,
				EndToEndID:       origEndToEndId,
				TransactionID:    origTxId,
				UETR:             origUETR,
				CreationDateTime: origCreDtTm,
			},
			CancellationReason: cxlReason,
			AdditionalInfo:     addtlInfo,
			SenderDI: FedNowDepositoryInstitution{
				SenderABANumber: senderABANumber,
			},
			ReceiverDI: FedNowDepositoryInstitution{
				ReceiverABANumber: receiverABANumber,
			},
		},
	}

	return &msg, nil
}

func extractClrSysMemberID(party head.Party44Choice) camt_056_001_08.Max35Text {
	if party.FIId == nil || party.FIId.FinInstnId.ClrSysMmbId == nil {
		return ""
	}
	return camt_056_001_08.Max35Text(party.FIId.FinInstnId.ClrSysMmbId.MmbId)
}
