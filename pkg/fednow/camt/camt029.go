package camt

import (
	"encoding/json"
	"encoding/xml"
	"time"

	camt_029_001_09 "github.com/mbanq/iso20022-go/ISO20022/camt_029_001_09"
	head "github.com/mbanq/iso20022-go/ISO20022/head_001_001_02"
	"github.com/mbanq/iso20022-go/pkg/common"
	"github.com/mbanq/iso20022-go/pkg/fednow/config"
)

func BuildCamt029Struct(message FedNowMessageCxlRsp, msgConfig *config.Config) (*camt_029_001_09.Document, error) {
	fedMsg := message.FedNowMsg

	clearingSystemId := camt_029_001_09.ExternalClearingSystemIdentification1Code(msgConfig.ClearingSystemId)

	assgnr := camt_029_001_09.Party40Choice{
		Agt: &camt_029_001_09.BranchAndFinancialInstitutionIdentification6{
			FinInstnId: camt_029_001_09.FinancialInstitutionIdentification18{
				ClrSysMmbId: &camt_029_001_09.ClearingSystemMemberIdentification2{
					MmbId: camt_029_001_09.Max35Text(fedMsg.SenderDI.SenderABANumber),
					ClrSysId: &camt_029_001_09.ClearingSystemIdentification2Choice{
						Cd: &clearingSystemId,
					},
				},
			},
		},
	}
	assgne := camt_029_001_09.Party40Choice{
		Agt: &camt_029_001_09.BranchAndFinancialInstitutionIdentification6{
			FinInstnId: camt_029_001_09.FinancialInstitutionIdentification18{
				ClrSysMmbId: &camt_029_001_09.ClearingSystemMemberIdentification2{
					MmbId: camt_029_001_09.Max35Text(fedMsg.ReceiverDI.ReceiverABANumber),
					ClrSysId: &camt_029_001_09.ClearingSystemIdentification2Choice{
						Cd: &clearingSystemId,
					},
				},
			},
		},
	}

	var resolvedCase *camt_029_001_09.Case5
	if fedMsg.ResolvedCase.CaseID != "" {
		creatorMemberID := diMemberID(fedMsg.ResolvedCase.CreatorDI)
		resolvedCase = &camt_029_001_09.Case5{
			Id: camt_029_001_09.Max35Text(fedMsg.ResolvedCase.CaseID),
			Cretr: camt_029_001_09.Party40Choice{
				Agt: &camt_029_001_09.BranchAndFinancialInstitutionIdentification6{
					FinInstnId: camt_029_001_09.FinancialInstitutionIdentification18{
						ClrSysMmbId: &camt_029_001_09.ClearingSystemMemberIdentification2{
							MmbId: camt_029_001_09.Max35Text(creatorMemberID),
							ClrSysId: &camt_029_001_09.ClearingSystemIdentification2Choice{
								Cd: &clearingSystemId,
							},
						},
					},
				},
			},
		}
	}

	status := camt_029_001_09.InvestigationStatus5Choice{
		Conf:           fedMsg.InvestigationStatus.Confirmation,
		RjctdMod:       fedMsg.InvestigationStatus.RejectedModification,
		AssgnmtCxlConf: fedMsg.InvestigationStatus.AssignmentCancellationConfirmed,
	}
	if fedMsg.InvestigationStatus.DuplicateOf != nil {
		creatorMemberID := diMemberID(fedMsg.InvestigationStatus.DuplicateOf.CreatorDI)
		status.DplctOf = &camt_029_001_09.Case5{
			Id: camt_029_001_09.Max35Text(fedMsg.InvestigationStatus.DuplicateOf.CaseID),
			Cretr: camt_029_001_09.Party40Choice{
				Agt: &camt_029_001_09.BranchAndFinancialInstitutionIdentification6{
					FinInstnId: camt_029_001_09.FinancialInstitutionIdentification18{
						ClrSysMmbId: &camt_029_001_09.ClearingSystemMemberIdentification2{
							MmbId: camt_029_001_09.Max35Text(creatorMemberID),
							ClrSysId: &camt_029_001_09.ClearingSystemIdentification2Choice{
								Cd: &clearingSystemId,
							},
						},
					},
				},
			},
		}
	}

	var cxlDetails []camt_029_001_09.UnderlyingTransaction22
	for _, detail := range fedMsg.CancellationDetails {
		txInf := camt_029_001_09.PaymentTransaction102{
			OrgnlInstrId:    detail.OriginalInstructionID,
			OrgnlEndToEndId: detail.OriginalEndToEndID,
			OrgnlUETR:       detail.OriginalUETR,
		}

		if detail.OriginalGroupInfo != nil {
			var orgnlCreationTime *common.ISODateTime
			if !time.Time(detail.OriginalGroupInfo.CreationDateTime).IsZero() {
				value := detail.OriginalGroupInfo.CreationDateTime
				orgnlCreationTime = &value
			}

			txInf.OrgnlGrpInf = &camt_029_001_09.OriginalGroupInformation29{
				OrgnlMsgId:   detail.OriginalGroupInfo.MessageID,
				OrgnlMsgNmId: detail.OriginalGroupInfo.MessageType,
				OrgnlCreDtTm: orgnlCreationTime,
			}
		}

		if detail.ResolutionRelatedInfo != nil {
			txInf.RsltnRltdInf = &camt_029_001_09.ResolutionData1{
				EndToEndId:     detail.ResolutionRelatedInfo.EndToEndID,
				TxId:           detail.ResolutionRelatedInfo.TransactionID,
				UETR:           detail.ResolutionRelatedInfo.UETR,
				IntrBkSttlmAmt: detail.ResolutionRelatedInfo.InterbankSettlementAmount,
				IntrBkSttlmDt:  detail.ResolutionRelatedInfo.InterbankSettlementDate,
			}
		}

		cxlDetails = append(cxlDetails, camt_029_001_09.UnderlyingTransaction22{
			TxInfAndSts: []camt_029_001_09.PaymentTransaction102{txInf},
		})
	}

	doc := &camt_029_001_09.Document{
		XMLName: xml.Name{Space: "urn:iso:std:iso:20022:tech:xsd:camt.029.001.09", Local: "Document"},
		RsltnOfInvstgtn: camt_029_001_09.ResolutionOfInvestigationV09{
			Assgnmt: camt_029_001_09.CaseAssignment5{
				Id:      camt_029_001_09.Max35Text(fedMsg.Identifier.MessageID),
				Assgnr:  assgnr,
				Assgne:  assgne,
				CreDtTm: fedMsg.CreationDateTime,
			},
			RslvdCase: resolvedCase,
			Sts:       status,
			CxlDtls:   cxlDetails,
		},
	}

	return doc, nil
}

func BuildCamt029(payload []byte, cfg *config.Config) (*camt_029_001_09.Document, error) {
	var message FedNowMessageCxlRsp
	if err := json.Unmarshal(payload, &message); err != nil {
		return nil, err
	}
	return BuildCamt029Struct(message, cfg)
}

func ParseCamt029(appHdr head.BusinessApplicationHeaderV02, document camt_029_001_09.Document) (*FedNowMessageCxlRsp, error) {
	response := document.RsltnOfInvstgtn

	var resolvedCase FedNowCase
	if response.RslvdCase != nil {
		resolvedCase = FedNowCase{
			CaseID: camt_029_001_09.Max35Text(response.RslvdCase.Id),
			CreatorDI: FedNowDepositoryInstitution2{
				SenderABANumber: extractCaseCreatorMemberID(response.RslvdCase.Cretr),
			},
		}
	}

	status := FedNowInvestigationStatus{
		Confirmation:                    response.Sts.Conf,
		RejectedModification:            response.Sts.RjctdMod,
		AssignmentCancellationConfirmed: response.Sts.AssgnmtCxlConf,
	}
	if response.Sts.DplctOf != nil {
		status.DuplicateOf = &FedNowCase{
			CaseID: camt_029_001_09.Max35Text(response.Sts.DplctOf.Id),
			CreatorDI: FedNowDepositoryInstitution2{
				SenderABANumber: extractCaseCreatorMemberID(response.Sts.DplctOf.Cretr),
			},
		}
	}

	var details []FedNowCxlRspDetails
	for _, cxl := range response.CxlDtls {
		for _, tx := range cxl.TxInfAndSts {
			detail := FedNowCxlRspDetails{
				OriginalInstructionID: tx.OrgnlInstrId,
				OriginalEndToEndID:    tx.OrgnlEndToEndId,
				OriginalUETR:          tx.OrgnlUETR,
			}

			if tx.OrgnlGrpInf != nil {
				var orgnlCreationTime common.ISODateTime
				if tx.OrgnlGrpInf.OrgnlCreDtTm != nil {
					orgnlCreationTime = common.ISODateTime(*tx.OrgnlGrpInf.OrgnlCreDtTm)
				}

				detail.OriginalGroupInfo = &FedNowOriginalGroupInfo{
					MessageID:        camt_029_001_09.Max35Text(tx.OrgnlGrpInf.OrgnlMsgId),
					MessageType:      camt_029_001_09.Max35Text(tx.OrgnlGrpInf.OrgnlMsgNmId),
					CreationDateTime: orgnlCreationTime,
				}
			}

			if tx.RsltnRltdInf != nil {
				detail.ResolutionRelatedInfo = &FedNowResolutionRelatedInfo{
					EndToEndID:                tx.RsltnRltdInf.EndToEndId,
					TransactionID:             tx.RsltnRltdInf.TxId,
					UETR:                      tx.RsltnRltdInf.UETR,
					InterbankSettlementAmount: tx.RsltnRltdInf.IntrBkSttlmAmt,
					InterbankSettlementDate:   tx.RsltnRltdInf.IntrBkSttlmDt,
				}
			}

			details = append(details, detail)
		}
	}

	senderABANumber := extractClrSysMemberIDCamt029(appHdr.Fr)
	receiverABANumber := extractClrSysMemberIDCamt029(appHdr.To)

	msg := FedNowMessageCxlRsp{
		FedNowMsg: FedNowCxlRsp{
			CreationDateTime: response.Assgnmt.CreDtTm,
			Identifier: FedNowIdentifierCxlRsp{
				BusinessMessageID: camt_029_001_09.Max35Text(appHdr.BizMsgIdr),
				MessageID:         camt_029_001_09.Max35Text(response.Assgnmt.Id),
				MessageType:       camt_029_001_09.Max35Text(appHdr.MsgDefIdr),
				CreationDateTime:  common.ISODateTime(appHdr.CreDt),
			},
			ResolvedCase:        resolvedCase,
			InvestigationStatus: status,
			CancellationDetails: details,
			SenderDI: FedNowDepositoryInstitution2{
				SenderABANumber: senderABANumber,
			},
			ReceiverDI: FedNowDepositoryInstitution2{
				ReceiverABANumber: receiverABANumber,
			},
		},
	}

	return &msg, nil
}

func extractClrSysMemberIDCamt029(party head.Party44Choice) camt_029_001_09.Max35Text {
	if party.FIId == nil || party.FIId.FinInstnId.ClrSysMmbId == nil {
		return ""
	}
	return camt_029_001_09.Max35Text(party.FIId.FinInstnId.ClrSysMmbId.MmbId)
}

func extractCaseCreatorMemberID(party camt_029_001_09.Party40Choice) camt_029_001_09.Max35Text {
	if party.Agt == nil || party.Agt.FinInstnId.ClrSysMmbId == nil {
		return ""
	}
	return camt_029_001_09.Max35Text(party.Agt.FinInstnId.ClrSysMmbId.MmbId)
}

func diMemberID(di FedNowDepositoryInstitution2) camt_029_001_09.Max35Text {
	if di.SenderABANumber != "" {
		return di.SenderABANumber
	}
	if di.ReceiverABANumber != "" {
		return di.ReceiverABANumber
	}
	return ""
}
