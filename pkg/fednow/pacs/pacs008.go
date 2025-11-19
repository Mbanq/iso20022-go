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

	if fedMsg.Identifier.UETR != "" {
		pacsDoc.FIToFICstmrCdtTrf.CdtTrfTxInf[0].PmtId.UETR = &fedMsg.Identifier.UETR
	}

	if *fedMsg.Identifier.TransactionID != "" {
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
			},
			PaymentType: FedNowPaymentType{
				CategoryPurpose: (*pacs_008_001_08.ExternalCategoryPurpose1Code)(cdtrftxinf.PmtTpInf.CtgyPurp.Prtry),
			},
			Amount: FedNowAmount{
				Text: json.Number(cdtrftxinf.IntrBkSttlmAmt.Text),
				Ccy:  cdtrftxinf.IntrBkSttlmAmt.Ccy,
			},
			SenderDI: FedNowDepositoryInstitution{
				ReceiverABANumber: pacs_008_001_08.Max35Text(appHdr.To.FIId.FinInstnId.ClrSysMmbId.MmbId),
			},
			ReceiverDI: FedNowDepositoryInstitution{
				SenderABANumber: pacs_008_001_08.Max35Text(appHdr.Fr.FIId.FinInstnId.ClrSysMmbId.MmbId),
			},
			Originator: FedNowParty{
				Personal: FedNowPersonal{
					Name: cdtrftxinf.Dbtr.Nm,
					Address: FedNowPstlAdr{
						StreetName:         cdtrftxinf.Dbtr.PstlAdr.StrtNm,
						BuildingNumber:     cdtrftxinf.Dbtr.PstlAdr.BldgNb,
						TownName:           cdtrftxinf.Dbtr.PstlAdr.TwnNm,
						CountrySubdivision: cdtrftxinf.Dbtr.PstlAdr.CtrySubDvsn,
						PostalCode:         cdtrftxinf.Dbtr.PstlAdr.PstCd,
						Country:            cdtrftxinf.Dbtr.PstlAdr.Ctry,
					},
					Identifier: pacs_008_001_08.Max34Text(cdtrftxinf.DbtrAcct.Id.Othr.Id),
				},
			},
			Beneficiary: FedNowParty{
				Personal: FedNowPersonal{
					Name: cdtrftxinf.Cdtr.Nm,
					Address: FedNowPstlAdr{
						StreetName:         cdtrftxinf.Cdtr.PstlAdr.StrtNm,
						BuildingNumber:     cdtrftxinf.Cdtr.PstlAdr.BldgNb,
						PostBox:            cdtrftxinf.Cdtr.PstlAdr.PstBx,
						TownName:           cdtrftxinf.Cdtr.PstlAdr.TwnNm,
						CountrySubdivision: cdtrftxinf.Cdtr.PstlAdr.CtrySubDvsn,
						PostalCode:         cdtrftxinf.Cdtr.PstlAdr.PstCd,
						Country:            cdtrftxinf.Cdtr.PstlAdr.Ctry,
					},
					Identifier: pacs_008_001_08.Max34Text(cdtrftxinf.CdtrAcct.Id.Othr.Id),
				},
			},
		},
	}

	return &fednowMsg, nil
}
