package pacs

import (
	"encoding/json"
	"encoding/xml"
	"fmt"

	head "github.com/mbanq/iso20022-go/ISO20022/head_001_001_02"
	"github.com/mbanq/iso20022-go/ISO20022/pacs_008_001_08"
	"github.com/mbanq/iso20022-go/pkg/common"
	"github.com/mbanq/iso20022-go/pkg/fednow/config"
)

func BuildPacs008Struct(message FedNowMessageCCT, msgConfig *config.Config) (*pacs_008_001_08.Document, error) {

	fedMsg := message.FedNowMsg

	// Assigning Configuration Values
	cd := pacs_008_001_08.ExternalCashClearingSystem1Code(msgConfig.ClearingSystem)
	clearingSystemId := pacs_008_001_08.ExternalClearingSystemIdentification1Code(msgConfig.ClearingSystemId)
	categoryPurpose := pacs_008_001_08.Max35Text(*fedMsg.PaymentType.CategoryPurpose)

	if fedMsg.Identifier.EndToEndID == "" {
		fedMsg.Identifier.EndToEndID = "NOTPROVIDED"
	}

	// Address Validation
	if err := fedMsg.Originator.Personal.Address.ValidateAddress(); err != nil {
		return nil, fmt.Errorf("invalid originator address: %w", err)
	}
	if err := fedMsg.Beneficiary.Personal.Address.ValidateAddress(); err != nil {
		return nil, fmt.Errorf("invalid beneficiary address: %w", err)
	}

	// Amount Validation
	amountFloat, err := fedMsg.Amount.Text.Float64()
	if err != nil {
		return nil, fmt.Errorf("invalid amount format: %w", err)
	}

	// Building the Pacs008 Struct
	pacsDoc := &pacs_008_001_08.Document{
		XMLName: xml.Name{Space: "urn:iso:std:iso:20022:tech:xsd:pacs.008.001.08", Local: "Document"},
		FIToFICstmrCdtTrf: pacs_008_001_08.FIToFICustomerCreditTransferV08{
			GrpHdr: pacs_008_001_08.GroupHeader93{
				MsgId:   pacs_008_001_08.Max35Text(fedMsg.Identifier.MessageID),
				CreDtTm: (common.ISODateTime)(fedMsg.CreationDateTime),
				NbOfTxs: "1",
				SttlmInf: pacs_008_001_08.SettlementInstruction7{
					SttlmMtd: msgConfig.SettlementMethod,
					ClrSys: &pacs_008_001_08.ClearingSystemIdentification3Choice{
						Cd: &cd,
					},
				},
			},
			CdtTrfTxInf: []pacs_008_001_08.CreditTransferTransaction39{
				{
					PmtId: pacs_008_001_08.PaymentIdentification7{
						InstrId:    fedMsg.Identifier.InstructionID,
						EndToEndId: fedMsg.Identifier.EndToEndID,
					},
					PmtTpInf: &pacs_008_001_08.PaymentTypeInformation28{
						LclInstrm: &msgConfig.LocalInstrument,
						CtgyPurp: &pacs_008_001_08.CategoryPurpose1Choice{
							Prtry: &categoryPurpose,
						},
					},
					IntrBkSttlmAmt: pacs_008_001_08.ActiveCurrencyAndAmount{
						Ccy:  fedMsg.Amount.Ccy,
						Text: fmt.Sprintf("%.2f", amountFloat),
					},
					IntrBkSttlmDt: (*common.ISODate)(&fedMsg.CreationDateTime),
					ChrgBr:        msgConfig.ChargeBearer,
					InstgAgt: &pacs_008_001_08.BranchAndFinancialInstitutionIdentification6{
						FinInstnId: pacs_008_001_08.FinancialInstitutionIdentification18{
							ClrSysMmbId: &pacs_008_001_08.ClearingSystemMemberIdentification2{
								MmbId: fedMsg.SenderDI.SenderABANumber,
								ClrSysId: &pacs_008_001_08.ClearingSystemIdentification2Choice{
									Cd: &clearingSystemId,
								},
							},
						},
					},
					InstdAgt: &pacs_008_001_08.BranchAndFinancialInstitutionIdentification6{
						FinInstnId: pacs_008_001_08.FinancialInstitutionIdentification18{
							ClrSysMmbId: &pacs_008_001_08.ClearingSystemMemberIdentification2{
								MmbId: fedMsg.ReceiverDI.ReceiverABANumber,
								ClrSysId: &pacs_008_001_08.ClearingSystemIdentification2Choice{
									Cd: &clearingSystemId,
								},
							},
						},
					},
					Dbtr: pacs_008_001_08.PartyIdentification135{
						Nm: fedMsg.Originator.Personal.Name,
						PstlAdr: &pacs_008_001_08.PostalAddress24{
							StrtNm:      fedMsg.Originator.Personal.Address.StreetName,
							BldgNb:      fedMsg.Originator.Personal.Address.BuildingNumber,
							TwnNm:       fedMsg.Originator.Personal.Address.TownName,
							CtrySubDvsn: fedMsg.Originator.Personal.Address.CountrySubdivision,
							PstCd:       fedMsg.Originator.Personal.Address.PostalCode,
							Ctry:        fedMsg.Originator.Personal.Address.Country,
						},
					},
					DbtrAcct: &pacs_008_001_08.CashAccount38{
						Id: pacs_008_001_08.AccountIdentification4Choice{
							Othr: &pacs_008_001_08.GenericAccountIdentification1{
								Id: fedMsg.Originator.Personal.Identifier,
							},
						},
					},
					DbtrAgt: pacs_008_001_08.BranchAndFinancialInstitutionIdentification6{
						FinInstnId: pacs_008_001_08.FinancialInstitutionIdentification18{
							ClrSysMmbId: &pacs_008_001_08.ClearingSystemMemberIdentification2{
								MmbId: fedMsg.SenderDI.SenderABANumber,
								ClrSysId: &pacs_008_001_08.ClearingSystemIdentification2Choice{
									Cd: &clearingSystemId,
								},
							},
							Nm: fedMsg.SenderDI.Name,
						},
					},
					CdtrAgt: pacs_008_001_08.BranchAndFinancialInstitutionIdentification6{
						FinInstnId: pacs_008_001_08.FinancialInstitutionIdentification18{
							ClrSysMmbId: &pacs_008_001_08.ClearingSystemMemberIdentification2{
								MmbId: fedMsg.ReceiverDI.ReceiverABANumber,
								ClrSysId: &pacs_008_001_08.ClearingSystemIdentification2Choice{
									Cd: &clearingSystemId,
								},
							},
						},
					},
					Cdtr: pacs_008_001_08.PartyIdentification135{
						Nm: fedMsg.Beneficiary.Personal.Name,
						PstlAdr: &pacs_008_001_08.PostalAddress24{
							StrtNm:      fedMsg.Beneficiary.Personal.Address.StreetName,
							BldgNb:      fedMsg.Beneficiary.Personal.Address.BuildingNumber,
							TwnNm:       fedMsg.Beneficiary.Personal.Address.TownName,
							CtrySubDvsn: fedMsg.Beneficiary.Personal.Address.CountrySubdivision,
							PstCd:       fedMsg.Beneficiary.Personal.Address.PostalCode,
							Ctry:        fedMsg.Beneficiary.Personal.Address.Country,
						},
					},
					CdtrAcct: &pacs_008_001_08.CashAccount38{
						Id: pacs_008_001_08.AccountIdentification4Choice{
							Othr: &pacs_008_001_08.GenericAccountIdentification1{
								Id: fedMsg.Beneficiary.Personal.Identifier,
							},
						},
					},
				},
			},
		},
	}

	if fedMsg.Identifier.UETR != nil {
		pacsDoc.FIToFICstmrCdtTrf.CdtTrfTxInf[0].PmtId.UETR = fedMsg.Identifier.UETR
	}

	if fedMsg.Identifier.TransactionID != nil && *fedMsg.Identifier.TransactionID != "" {
		pacsDoc.FIToFICstmrCdtTrf.CdtTrfTxInf[0].PmtId.TxId = fedMsg.Identifier.TransactionID
	}
	return pacsDoc, nil

}

func BuildPacs008(payload []byte, config *config.Config) (*pacs_008_001_08.Document, error) {

	var message FedNowMessageCCT
	if err := json.Unmarshal(payload, &message); err != nil {
		return nil, err
	}

	return BuildPacs008Struct(message, config)
}

func ParsePacs008(appHdr head.BusinessApplicationHeaderV02, document pacs_008_001_08.Document) (*FedNowMessageCCT, error) {

	fitoficstmrcdttrf := document.FIToFICstmrCdtTrf
	cdtrftxinf := fitoficstmrcdttrf.CdtTrfTxInf[0]

	categoryPurpose := resolveCategoryPurpose(cdtrftxinf.PmtTpInf)
	originatorAddress := convertPostalAddress(cdtrftxinf.Dbtr.PstlAdr)
	beneficiaryAddress := convertPostalAddress(cdtrftxinf.Cdtr.PstlAdr)
	originatorIdentifier := extractAccountIdentifier(cdtrftxinf.DbtrAcct)
	beneficiaryIdentifier := extractAccountIdentifier(cdtrftxinf.CdtrAcct)
	senderABANumber := extractClrSysMemberIDFromAgent(cdtrftxinf.InstgAgt)
	if senderABANumber == "" {
		senderABANumber = extractClrSysMemberID(appHdr.To)
	}
	receiverABANumber := extractClrSysMemberIDFromAgent(cdtrftxinf.InstdAgt)
	if receiverABANumber == "" {
		receiverABANumber = extractClrSysMemberID(appHdr.Fr)
	}

	var uetr pacs_008_001_08.UUIDv4Identifier
	if cdtrftxinf.PmtId.UETR != nil {
		uetr = pacs_008_001_08.UUIDv4Identifier(*cdtrftxinf.PmtId.UETR)
	}

	fednowMsg := FedNowMessageCCT{
		FedNowMsg: FedNowDetails{
			CreationDateTime: common.ISODateTime(fitoficstmrcdttrf.GrpHdr.CreDtTm),
			Identifier: FedNowIdentifier{
				BusinessMessageID: pacs_008_001_08.Max35Text(appHdr.BizMsgIdr),
				MessageID:         pacs_008_001_08.Max35Text(fitoficstmrcdttrf.GrpHdr.MsgId),
				InstructionID:     cdtrftxinf.PmtId.InstrId,
				EndToEndID:        cdtrftxinf.PmtId.EndToEndId,
				TransactionID:     cdtrftxinf.PmtId.TxId,
				CreationDateTime:  common.ISODateTime(appHdr.CreDt),
				UETR:              uetr,
			},
			PaymentType: FedNowPaymentType{
				CategoryPurpose: categoryPurpose,
			},
			Amount: FedNowAmount{
				Text: json.Number(cdtrftxinf.IntrBkSttlmAmt.Text),
				Ccy:  cdtrftxinf.IntrBkSttlmAmt.Ccy,
			},
			SenderDI: FedNowDepositoryInstitution{
				SenderABANumber: senderABANumber,
			},
			ReceiverDI: FedNowDepositoryInstitution{
				ReceiverABANumber: receiverABANumber,
			},
			Originator: FedNowParty{
				Personal: FedNowPersonal{
					Name:       cdtrftxinf.Dbtr.Nm,
					Address:    originatorAddress,
					Identifier: originatorIdentifier,
				},
			},
			Beneficiary: FedNowParty{
				Personal: FedNowPersonal{
					Name:       cdtrftxinf.Cdtr.Nm,
					Address:    beneficiaryAddress,
					Identifier: beneficiaryIdentifier,
				},
			},
		},
	}
	// Map inbound pacs.008 UETR (optional)
	fednowMsg.FedNowMsg.Identifier.UETR = cdtrftxinf.PmtId.UETR

	return &fednowMsg, nil
}

func resolveCategoryPurpose(pmtType *pacs_008_001_08.PaymentTypeInformation28) *pacs_008_001_08.ExternalCategoryPurpose1Code {
	if pmtType == nil || pmtType.CtgyPurp == nil {
		return nil
	}

	if pmtType.CtgyPurp.Cd != nil {
		return pmtType.CtgyPurp.Cd
	}

	if pmtType.CtgyPurp.Prtry != nil {
		cp := pacs_008_001_08.ExternalCategoryPurpose1Code(*pmtType.CtgyPurp.Prtry)
		return &cp
	}

	return nil
}

func convertPostalAddress(addr *pacs_008_001_08.PostalAddress24) FedNowPstlAdr {
	if addr == nil {
		return FedNowPstlAdr{}
	}

	return FedNowPstlAdr{
		StreetName:         addr.StrtNm,
		BuildingNumber:     addr.BldgNb,
		PostBox:            addr.PstBx,
		TownName:           addr.TwnNm,
		CountrySubdivision: addr.CtrySubDvsn,
		PostalCode:         addr.PstCd,
		Country:            addr.Ctry,
	}
}

func extractAccountIdentifier(acct *pacs_008_001_08.CashAccount38) pacs_008_001_08.Max34Text {
	if acct == nil {
		return ""
	}

	if acct.Id.Othr != nil {
		return acct.Id.Othr.Id
	}

	if acct.Id.IBAN != nil {
		return pacs_008_001_08.Max34Text(*acct.Id.IBAN)
	}

	return ""
}

func extractClrSysMemberID(party head.Party44Choice) pacs_008_001_08.Max35Text {
	if party.FIId == nil || party.FIId.FinInstnId.ClrSysMmbId == nil {
		return ""
	}

	return pacs_008_001_08.Max35Text(party.FIId.FinInstnId.ClrSysMmbId.MmbId)
}

func extractClrSysMemberIDFromAgent(agent *pacs_008_001_08.BranchAndFinancialInstitutionIdentification6) pacs_008_001_08.Max35Text {
	if agent == nil {
		return ""
	}
	if agent.FinInstnId.ClrSysMmbId == nil {
		return ""
	}
	return agent.FinInstnId.ClrSysMmbId.MmbId
}
