package camt

import (
	"encoding/json"
	"encoding/xml"
	"time"

	camt_026_001_07 "github.com/mbanq/iso20022-go/ISO20022/camt_026_001_07"
	head "github.com/mbanq/iso20022-go/ISO20022/head_001_001_02"
	"github.com/mbanq/iso20022-go/pkg/common"
	"github.com/mbanq/iso20022-go/pkg/fednow/config"
)

func BuildCamt026Struct(message FedNowMessageInfoReq, msgConfig *config.Config) (*camt_026_001_07.Document, error) {
	fedMsg := message.FedNowMsg

	clearingSystemId := camt_026_001_07.ExternalClearingSystemIdentification1Code(msgConfig.ClearingSystemId)

	assgnr := camt_026_001_07.Party40Choice{
		Agt: &camt_026_001_07.BranchAndFinancialInstitutionIdentification6{
			FinInstnId: camt_026_001_07.FinancialInstitutionIdentification18{
				ClrSysMmbId: &camt_026_001_07.ClearingSystemMemberIdentification2{
					MmbId: camt_026_001_07.Max35Text(fedMsg.SenderDI.SenderABANumber),
					ClrSysId: &camt_026_001_07.ClearingSystemIdentification2Choice{
						Cd: &clearingSystemId,
					},
				},
			},
		},
	}
	assgne := camt_026_001_07.Party40Choice{
		Agt: &camt_026_001_07.BranchAndFinancialInstitutionIdentification6{
			FinInstnId: camt_026_001_07.FinancialInstitutionIdentification18{
				ClrSysMmbId: &camt_026_001_07.ClearingSystemMemberIdentification2{
					MmbId: camt_026_001_07.Max35Text(fedMsg.ReceiverDI.ReceiverABANumber),
					ClrSysId: &camt_026_001_07.ClearingSystemIdentification2Choice{
						Cd: &clearingSystemId,
					},
				},
			},
		},
	}

	// Case
	var caseElement *camt_026_001_07.Case5
	if fedMsg.Case.CaseID != "" {
		creatorMemberID := infoReqDIMemberID(fedMsg.Case.CreatorDI)
		caseElement = &camt_026_001_07.Case5{
			Id: camt_026_001_07.Max35Text(fedMsg.Case.CaseID),
			Cretr: camt_026_001_07.Party40Choice{
				Agt: &camt_026_001_07.BranchAndFinancialInstitutionIdentification6{
					FinInstnId: camt_026_001_07.FinancialInstitutionIdentification18{
						ClrSysMmbId: &camt_026_001_07.ClearingSystemMemberIdentification2{
							MmbId: camt_026_001_07.Max35Text(creatorMemberID),
							ClrSysId: &camt_026_001_07.ClearingSystemIdentification2Choice{
								Cd: &clearingSystemId,
							},
						},
					},
				},
			},
		}
	}

	// Underlying
	underlying := buildUnderlyingChoice(fedMsg.Underlying)

	// Justification
	justification := buildJustification(fedMsg.Justification)

	doc := &camt_026_001_07.Document{
		XMLName: xml.Name{Space: "urn:iso:std:iso:20022:tech:xsd:camt.026.001.07", Local: "Document"},
		UblToApply: camt_026_001_07.UnableToApplyV07{
			Assgnmt: camt_026_001_07.CaseAssignment5{
				Id:      camt_026_001_07.Max35Text(fedMsg.Identifier.MessageID),
				Assgnr:  assgnr,
				Assgne:  assgne,
				CreDtTm: fedMsg.CreationDateTime,
			},
			Case:   caseElement,
			Undrlyg: underlying,
			Justfn: justification,
		},
	}

	return doc, nil
}

func buildUnderlyingChoice(u FedNowUnderlyingInfoReq) camt_026_001_07.UnderlyingTransaction5Choice {
	var choice camt_026_001_07.UnderlyingTransaction5Choice

	if u.Initiation != nil {
		initn := u.Initiation
		pi := &camt_026_001_07.UnderlyingPaymentInstruction5{
			OrgnlPmtInfId:   initn.OriginalPaymentInfoID,
			OrgnlInstrId:    initn.OriginalInstructionID,
			OrgnlEndToEndId: initn.OriginalEndToEndID,
			OrgnlUETR:       initn.OriginalUETR,
		}

		if initn.OriginalInstructedAmount != nil {
			pi.OrgnlInstdAmt = *initn.OriginalInstructedAmount
		}

		if initn.OriginalGroupInfo != nil {
			var orgnlCreDtTm *common.ISODateTime
			if !time.Time(initn.OriginalGroupInfo.CreationDateTime).IsZero() {
				v := initn.OriginalGroupInfo.CreationDateTime
				orgnlCreDtTm = &v
			}
			pi.OrgnlGrpInf = &camt_026_001_07.UnderlyingGroupInformation1{
				OrgnlMsgId:   camt_026_001_07.Max35Text(initn.OriginalGroupInfo.MessageID),
				OrgnlMsgNmId: camt_026_001_07.Max35Text(initn.OriginalGroupInfo.MessageType),
				OrgnlCreDtTm: orgnlCreDtTm,
			}
		}

		if initn.RequestedExecutionDate != nil {
			pi.ReqdExctnDt = &camt_026_001_07.DateAndDateTime2Choice{
				Dt:   initn.RequestedExecutionDate.Date,
				DtTm: initn.RequestedExecutionDate.DateTime,
			}
		}

		choice.Initn = pi
	}

	if u.Interbank != nil {
		ib := u.Interbank
		pt := &camt_026_001_07.UnderlyingPaymentTransaction4{
			OrgnlInstrId:    ib.OriginalInstructionID,
			OrgnlEndToEndId: ib.OriginalEndToEndID,
			OrgnlTxId:       ib.OriginalTransactionID,
			OrgnlUETR:       ib.OriginalUETR,
		}

		if ib.OriginalInterbankSettlementAmount != nil {
			pt.OrgnlIntrBkSttlmAmt = *ib.OriginalInterbankSettlementAmount
		}

		if ib.OriginalInterbankSettlementDate != nil {
			pt.OrgnlIntrBkSttlmDt = *ib.OriginalInterbankSettlementDate
		}

		if ib.OriginalGroupInfo != nil {
			var orgnlCreDtTm *common.ISODateTime
			if !time.Time(ib.OriginalGroupInfo.CreationDateTime).IsZero() {
				v := ib.OriginalGroupInfo.CreationDateTime
				orgnlCreDtTm = &v
			}
			pt.OrgnlGrpInf = &camt_026_001_07.UnderlyingGroupInformation1{
				OrgnlMsgId:   camt_026_001_07.Max35Text(ib.OriginalGroupInfo.MessageID),
				OrgnlMsgNmId: camt_026_001_07.Max35Text(ib.OriginalGroupInfo.MessageType),
				OrgnlCreDtTm: orgnlCreDtTm,
			}
		}

		choice.IntrBk = pt
	}

	return choice
}

func buildJustification(j FedNowJustification) camt_026_001_07.UnableToApplyJustification3Choice {
	var choice camt_026_001_07.UnableToApplyJustification3Choice

	if j.MissingOrIncorrectInfo != nil {
		moi := &camt_026_001_07.MissingOrIncorrectInformation3{
			AMLReq: j.MissingOrIncorrectInfo.AMLRequest,
		}

		for _, m := range j.MissingOrIncorrectInfo.MissingInfo {
			moi.MssngInf = append(moi.MssngInf, camt_026_001_07.UnableToApplyMissing1{
				Cd:            m.Code,
				AddtlMssngInf: m.AdditionalInfo,
			})
		}

		for _, i := range j.MissingOrIncorrectInfo.IncorrectInfo {
			moi.IncrrctInf = append(moi.IncrrctInf, camt_026_001_07.UnableToApplyIncorrect1{
				Cd:              i.Code,
				AddtlIncrrctInf: i.AdditionalInfo,
			})
		}

		choice.MssngOrIncrrctInf = moi
	}

	if j.PossibleDuplicateInstruction != nil {
		choice.PssblDplctInstr = j.PossibleDuplicateInstruction
	}

	return choice
}

func BuildCamt026(payload []byte, cfg *config.Config) (*camt_026_001_07.Document, error) {
	var message FedNowMessageInfoReq
	if err := json.Unmarshal(payload, &message); err != nil {
		return nil, err
	}
	return BuildCamt026Struct(message, cfg)
}

func ParseCamt026(appHdr head.BusinessApplicationHeaderV02, document camt_026_001_07.Document) (*FedNowMessageInfoReq, error) {
	uta := document.UblToApply

	// Case
	var caseInfo FedNowCaseInfoReq
	if uta.Case != nil {
		caseInfo = FedNowCaseInfoReq{
			CaseID: camt_026_001_07.Max35Text(uta.Case.Id),
			CreatorDI: FedNowDepositoryInstitutionInfo{
				SenderABANumber: extractCaseCreatorMemberIDCamt026(uta.Case.Cretr),
			},
		}
	}

	// Underlying
	underlying := parseUnderlying(uta.Undrlyg)

	// Justification
	justification := parseJustification(uta.Justfn)

	senderABANumber := extractClrSysMemberIDFromParty40ChoiceCamt026(uta.Assgnmt.Assgnr)
	if senderABANumber == "" {
		senderABANumber = extractClrSysMemberIDCamt026(appHdr.To)
	}
	receiverABANumber := extractClrSysMemberIDFromParty40ChoiceCamt026(uta.Assgnmt.Assgne)
	if receiverABANumber == "" {
		receiverABANumber = extractClrSysMemberIDCamt026(appHdr.Fr)
	}

	msg := FedNowMessageInfoReq{
		FedNowMsg: FedNowInfoReq{
			CreationDateTime: uta.Assgnmt.CreDtTm,
			Identifier: FedNowIdentifierInfoReq{
				BusinessMessageID: camt_026_001_07.Max35Text(appHdr.BizMsgIdr),
				MessageID:         camt_026_001_07.Max35Text(uta.Assgnmt.Id),
				MessageType:       camt_026_001_07.Max35Text(appHdr.MsgDefIdr),
				CreationDateTime:  common.ISODateTime(appHdr.CreDt),
			},
			Case:          caseInfo,
			Underlying:    underlying,
			Justification: justification,
			SenderDI: FedNowDepositoryInstitutionInfo{
				SenderABANumber: senderABANumber,
			},
			ReceiverDI: FedNowDepositoryInstitutionInfo{
				ReceiverABANumber: receiverABANumber,
			},
		},
	}

	return &msg, nil
}

func parseUnderlying(u camt_026_001_07.UnderlyingTransaction5Choice) FedNowUnderlyingInfoReq {
	var result FedNowUnderlyingInfoReq

	if u.Initn != nil {
		initn := &FedNowUnderlyingInitiation{
			OriginalPaymentInfoID: u.Initn.OrgnlPmtInfId,
			OriginalInstructionID: u.Initn.OrgnlInstrId,
			OriginalEndToEndID:    u.Initn.OrgnlEndToEndId,
			OriginalUETR:          u.Initn.OrgnlUETR,
			OriginalInstructedAmount: &u.Initn.OrgnlInstdAmt,
		}

		if u.Initn.OrgnlGrpInf != nil {
			var orgnlCreDtTm common.ISODateTime
			if u.Initn.OrgnlGrpInf.OrgnlCreDtTm != nil {
				orgnlCreDtTm = *u.Initn.OrgnlGrpInf.OrgnlCreDtTm
			}
			initn.OriginalGroupInfo = &FedNowOriginalGroupInfoCamt026{
				MessageID:        u.Initn.OrgnlGrpInf.OrgnlMsgId,
				MessageType:      u.Initn.OrgnlGrpInf.OrgnlMsgNmId,
				CreationDateTime: orgnlCreDtTm,
			}
		}

		if u.Initn.ReqdExctnDt != nil {
			initn.RequestedExecutionDate = &FedNowDateChoice{
				Date:     u.Initn.ReqdExctnDt.Dt,
				DateTime: u.Initn.ReqdExctnDt.DtTm,
			}
		}

		result.Initiation = initn
	}

	if u.IntrBk != nil {
		ib := &FedNowUnderlyingInterbank{
			OriginalInstructionID:             u.IntrBk.OrgnlInstrId,
			OriginalEndToEndID:                u.IntrBk.OrgnlEndToEndId,
			OriginalTransactionID:             u.IntrBk.OrgnlTxId,
			OriginalUETR:                      u.IntrBk.OrgnlUETR,
			OriginalInterbankSettlementAmount: &u.IntrBk.OrgnlIntrBkSttlmAmt,
			OriginalInterbankSettlementDate:   &u.IntrBk.OrgnlIntrBkSttlmDt,
		}

		if u.IntrBk.OrgnlGrpInf != nil {
			var orgnlCreDtTm common.ISODateTime
			if u.IntrBk.OrgnlGrpInf.OrgnlCreDtTm != nil {
				orgnlCreDtTm = *u.IntrBk.OrgnlGrpInf.OrgnlCreDtTm
			}
			ib.OriginalGroupInfo = &FedNowOriginalGroupInfoCamt026{
				MessageID:        u.IntrBk.OrgnlGrpInf.OrgnlMsgId,
				MessageType:      u.IntrBk.OrgnlGrpInf.OrgnlMsgNmId,
				CreationDateTime: orgnlCreDtTm,
			}
		}

		result.Interbank = ib
	}

	return result
}

func parseJustification(j camt_026_001_07.UnableToApplyJustification3Choice) FedNowJustification {
	var result FedNowJustification

	if j.MssngOrIncrrctInf != nil {
		moi := &FedNowMissingOrIncorrectInfo{
			AMLRequest: j.MssngOrIncrrctInf.AMLReq,
		}

		for _, m := range j.MssngOrIncrrctInf.MssngInf {
			moi.MissingInfo = append(moi.MissingInfo, FedNowMissingInfo{
				Code:           m.Cd,
				AdditionalInfo: m.AddtlMssngInf,
			})
		}

		for _, i := range j.MssngOrIncrrctInf.IncrrctInf {
			moi.IncorrectInfo = append(moi.IncorrectInfo, FedNowIncorrectInfo{
				Code:           i.Cd,
				AdditionalInfo: i.AddtlIncrrctInf,
			})
		}

		result.MissingOrIncorrectInfo = moi
	}

	if j.PssblDplctInstr != nil {
		result.PossibleDuplicateInstruction = j.PssblDplctInstr
	}

	return result
}

func extractClrSysMemberIDCamt026(party head.Party44Choice) camt_026_001_07.Max35Text {
	if party.FIId == nil || party.FIId.FinInstnId.ClrSysMmbId == nil {
		return ""
	}
	return camt_026_001_07.Max35Text(party.FIId.FinInstnId.ClrSysMmbId.MmbId)
}

func extractCaseCreatorMemberIDCamt026(party camt_026_001_07.Party40Choice) camt_026_001_07.Max35Text {
	if party.Agt == nil || party.Agt.FinInstnId.ClrSysMmbId == nil {
		return ""
	}
	return camt_026_001_07.Max35Text(party.Agt.FinInstnId.ClrSysMmbId.MmbId)
}

func extractClrSysMemberIDFromParty40ChoiceCamt026(party camt_026_001_07.Party40Choice) camt_026_001_07.Max35Text {
	if party.Agt == nil || party.Agt.FinInstnId.ClrSysMmbId == nil {
		return ""
	}
	return camt_026_001_07.Max35Text(party.Agt.FinInstnId.ClrSysMmbId.MmbId)
}

func infoReqDIMemberID(di FedNowDepositoryInstitutionInfo) camt_026_001_07.Max35Text {
	if di.SenderABANumber != "" {
		return di.SenderABANumber
	}
	if di.ReceiverABANumber != "" {
		return di.ReceiverABANumber
	}
	return ""
}
