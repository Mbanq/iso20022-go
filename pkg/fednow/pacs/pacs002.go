package pacs

import (
	"encoding/json"
	"encoding/xml"
	"time"

	head "github.com/mbanq/iso20022-go/ISO20022/head_001_001_02"
	"github.com/mbanq/iso20022-go/ISO20022/pacs_002_001_10"
	"github.com/mbanq/iso20022-go/ISO20022/pacs_008_001_08"
	"github.com/mbanq/iso20022-go/pkg/common"
	"github.com/mbanq/iso20022-go/pkg/fednow/config"
)

func BuildPacs002Struct(message FedNowMessageACK, msgConfig *config.Config) (*pacs_002_001_10.Document, error) {

	fedMsg := message.FedNowMsg

	clearingSystemId := pacs_002_001_10.ExternalClearingSystemIdentification1Code(msgConfig.ClearingSystemId)
	creationTime := common.ISODateTime(fedMsg.CreationDateTime)
	instructionId := pacs_002_001_10.Max35Text(*fedMsg.OriginalIdentifier.InstructionID)
	endToEndId := pacs_002_001_10.Max35Text(fedMsg.OriginalIdentifier.EndToEndID)
	uetr := pacs_002_001_10.UUIDv4Identifier(fedMsg.OriginalIdentifier.UETR)

	pacsDoc := &pacs_002_001_10.Document{
		XMLName: xml.Name{Space: "urn:iso:std:iso:20022:tech:xsd:pacs.002.001.10", Local: "Document"},
		FIToFIPmtStsRpt: pacs_002_001_10.FIToFIPaymentStatusReportV10{
			GrpHdr: pacs_002_001_10.GroupHeader91{
				MsgId:   pacs_002_001_10.Max35Text(fedMsg.Identifier.MessageID),
				CreDtTm: (common.ISODateTime)(fedMsg.CreationDateTime),
			},
			TxInfAndSts: []pacs_002_001_10.PaymentTransaction110{
				{
					OrgnlGrpInf: &pacs_002_001_10.OriginalGroupInformation29{
						OrgnlMsgId:   pacs_002_001_10.Max35Text(fedMsg.Identifier.MessageID),
						OrgnlMsgNmId: pacs_002_001_10.Max35Text(fedMsg.Identifier.MessageType),
						OrgnlCreDtTm: &creationTime,
					},
					OrgnlInstrId:    &instructionId,
					OrgnlEndToEndId: &endToEndId,
					OrgnlUETR:       &uetr,
					TxSts:           fedMsg.PaymentStatus.PaymentStatus,
					InstgAgt: &pacs_002_001_10.BranchAndFinancialInstitutionIdentification6{
						FinInstnId: pacs_002_001_10.FinancialInstitutionIdentification18{
							ClrSysMmbId: &pacs_002_001_10.ClearingSystemMemberIdentification2{
								MmbId: pacs_002_001_10.Max35Text(fedMsg.SenderDI.SenderABANumber),
								ClrSysId: &pacs_002_001_10.ClearingSystemIdentification2Choice{
									Cd: &clearingSystemId,
								},
							},
						},
					},
					InstdAgt: &pacs_002_001_10.BranchAndFinancialInstitutionIdentification6{
						FinInstnId: pacs_002_001_10.FinancialInstitutionIdentification18{
							ClrSysMmbId: &pacs_002_001_10.ClearingSystemMemberIdentification2{
								MmbId: pacs_002_001_10.Max35Text(fedMsg.ReceiverDI.ReceiverABANumber),
								ClrSysId: &pacs_002_001_10.ClearingSystemIdentification2Choice{
									Cd: &clearingSystemId,
								},
							},
						},
					},
				},
			},
		},
	}
	if *fedMsg.PaymentStatus.PaymentStatus == "ACSC" || *fedMsg.PaymentStatus.PaymentStatus == "ACWP" {
		if fedMsg.PaymentStatus.AcceptanceDateTime != nil {
			pacsDoc.FIToFIPmtStsRpt.TxInfAndSts[0].AccptncDtTm = fedMsg.PaymentStatus.AcceptanceDateTime
			acceptanceDate := common.ISODate(time.Time(*fedMsg.PaymentStatus.AcceptanceDateTime))
			pacsDoc.FIToFIPmtStsRpt.TxInfAndSts[0].FctvIntrBkSttlmDt = &pacs_002_001_10.DateAndDateTime2Choice{
				Dt: &acceptanceDate,
			}
		}
	}

	if *fedMsg.PaymentStatus.PaymentStatus == "RJCT" {
		stsRsnCode := pacs_002_001_10.ExternalStatusReason1Code(*fedMsg.PaymentStatus.StatusReason)
		pacsDoc.FIToFIPmtStsRpt.TxInfAndSts[0].StsRsnInf = []pacs_002_001_10.StatusReasonInformation12{
			{
				Rsn: &pacs_002_001_10.StatusReason6Choice{
					Cd: &stsRsnCode,
				},
				AddtlInf: []pacs_002_001_10.Max105Text{*fedMsg.PaymentStatus.AdditionalInformation},
			},
		}
	}

	return pacsDoc, nil
}

func BuildPacs002(payload []byte, config *config.Config) (*pacs_002_001_10.Document, error) {

	var message FedNowMessageACK
	if err := json.Unmarshal(payload, &message); err != nil {
		return nil, err
	}

	return BuildPacs002Struct(message, config)
}

func ParsePacs002(appHdr head.BusinessApplicationHeaderV02, document pacs_002_001_10.Document) (*FedNowMessageACK, error) {

	fitofipmtstsrpt := document.FIToFIPmtStsRpt
	txinfandsts := fitofipmtstsrpt.TxInfAndSts[0]

	var orgnlInstrId *pacs_008_001_08.Max35Text
	if txinfandsts.OrgnlInstrId != nil {
		val := pacs_008_001_08.Max35Text(*txinfandsts.OrgnlInstrId)
		orgnlInstrId = &val
	}

	var orgnlMsgId pacs_008_001_08.Max35Text
	var orgnlMsgNmId pacs_008_001_08.Max35Text
	var orgnlCreDtTm common.ISODateTime
	if txinfandsts.OrgnlGrpInf != nil {
		orgnlMsgId = pacs_008_001_08.Max35Text(txinfandsts.OrgnlGrpInf.OrgnlMsgId)
		orgnlMsgNmId = pacs_008_001_08.Max35Text(txinfandsts.OrgnlGrpInf.OrgnlMsgNmId)
		if txinfandsts.OrgnlGrpInf.OrgnlCreDtTm != nil {
			orgnlCreDtTm = common.ISODateTime(*txinfandsts.OrgnlGrpInf.OrgnlCreDtTm)
		}
	}

	var orgnlEndToEndId pacs_008_001_08.Max35Text
	if txinfandsts.OrgnlEndToEndId != nil {
		orgnlEndToEndId = pacs_008_001_08.Max35Text(*txinfandsts.OrgnlEndToEndId)
	}

	senderABANumber := extractClrSysMemberID(appHdr.Fr)
	receiverABANumber := extractClrSysMemberID(appHdr.To)

	fednowMsg := FedNowMessageACK{
		FedNowMsg: FedNowACK{
			CreationDateTime: common.ISODateTime(fitofipmtstsrpt.GrpHdr.CreDtTm),
			Identifier: FedNowIdentifier{
				BusinessMessageID: pacs_008_001_08.Max35Text(appHdr.BizMsgIdr),
				MessageID:         pacs_008_001_08.Max35Text(fitofipmtstsrpt.GrpHdr.MsgId),
				CreationDateTime:  common.ISODateTime(appHdr.CreDt),
			},
			OriginalIdentifier: FedNowIdentifier{
				MessageID:        orgnlMsgId,
				MessageType:      orgnlMsgNmId,
				InstructionID:    orgnlInstrId,
				EndToEndID:       orgnlEndToEndId,
				CreationDateTime: orgnlCreDtTm,
			},
			PaymentStatus: PaymentStatus{
				PaymentStatus: txinfandsts.TxSts,
			},
			SenderDI: FedNowDepositoryInstitution{
				SenderABANumber: senderABANumber,
			},
			ReceiverDI: FedNowDepositoryInstitution{
				ReceiverABANumber: receiverABANumber,
			},
		},
	}

	if txinfandsts.TxSts != nil {
		switch *txinfandsts.TxSts {
		case "ACSC", "ACWP":
			if txinfandsts.AccptncDtTm != nil {
				fednowMsg.FedNowMsg.PaymentStatus.AcceptanceDateTime = txinfandsts.AccptncDtTm
			}
		case "RJCT":
			if len(txinfandsts.StsRsnInf) > 0 {
				statusReasonInfo := txinfandsts.StsRsnInf[0]
				if statusReasonInfo.Rsn != nil && statusReasonInfo.Rsn.Cd != nil {
					reasonCode := pacs_002_001_10.ExternalStatusReason1Code(*statusReasonInfo.Rsn.Cd)
					fednowMsg.FedNowMsg.PaymentStatus.StatusReason = &reasonCode
				}
				if len(statusReasonInfo.AddtlInf) > 0 {
					additionalInfo := pacs_002_001_10.Max105Text(statusReasonInfo.AddtlInf[0])
					fednowMsg.FedNowMsg.PaymentStatus.AdditionalInformation = &additionalInfo
				}
			}
		}
	}

	return &fednowMsg, nil
}
