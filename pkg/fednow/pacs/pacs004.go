// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package pacs

import (
	pacs004 "github.com/mbanq/iso20022-go/ISO20022/pacs_004_001_10"
	"github.com/mbanq/iso20022-go/pkg/common"
	"github.com/mbanq/iso20022-go/pkg/fednow/config"
)

// BuildPacs004Struct creates a pacs.004.001.10 message from a FedNowMessageRtn struct
func BuildPacs004Struct(message FedNowMessageRtn, msgConfig *config.Config) (*pacs004.Document, error) {

	clearingSystemCd := pacs004.ExternalCashClearingSystem1Code(msgConfig.ClearingSystem)
	clearingSystemId := pacs004.ExternalClearingSystemIdentification1Code(msgConfig.ClearingSystemId)
	chargebearer := pacs004.ChargeBearerType1Code(msgConfig.ChargeBearer)
	OrgnlCreDtTm := message.FedNowMsg.OriginalIdentifier.CreationDateTime
	localInstrument := pacs004.Max35Text(*msgConfig.LocalInstrument.Prtry)
	var orgnlUetr *pacs004.UUIDv4Identifier
	if message.FedNowMsg.OriginalIdentifier.UETR != nil {
		tmp := pacs004.UUIDv4Identifier(*message.FedNowMsg.OriginalIdentifier.UETR)
		orgnlUetr = &tmp
	}

	pacsDoc := &pacs004.Document{
		PmtRtr: pacs004.PaymentReturnV10{
			GrpHdr: pacs004.GroupHeader90{
				MsgId:   pacs004.Max35Text(message.FedNowMsg.Identifier.MessageID),
				CreDtTm: message.FedNowMsg.CreationDateTime,
				NbOfTxs: "1",
				SttlmInf: pacs004.SettlementInstruction7{
					SttlmMtd: pacs004.SettlementMethod1Code(msgConfig.SettlementMethod),
					ClrSys: &pacs004.ClearingSystemIdentification3Choice{
						Cd: &clearingSystemCd,
					},
				},
			},
			TxInf: []pacs004.PaymentTransaction118{
				{
					OrgnlGrpInf: &pacs004.OriginalGroupInformation29{
						OrgnlMsgId:   pacs004.Max35Text(message.FedNowMsg.OriginalIdentifier.MessageID),
						OrgnlMsgNmId: pacs004.Max35Text(message.FedNowMsg.OriginalIdentifier.MessageType),
						OrgnlCreDtTm: &OrgnlCreDtTm,
					},
					OrgnlInstrId:    (*pacs004.Max35Text)(message.FedNowMsg.OriginalIdentifier.InstructionID),
					OrgnlEndToEndId: (*pacs004.Max35Text)(&message.FedNowMsg.OriginalIdentifier.EndToEndID),
					OrgnlUETR:       orgnlUetr,
					OrgnlIntrBkSttlmAmt: &pacs004.ActiveOrHistoricCurrencyAndAmount{
						Ccy:  pacs004.ActiveOrHistoricCurrencyCode(message.FedNowMsg.Amount.Ccy),
						Text: string(message.FedNowMsg.Amount.Text),
					},
					OrgnlIntrBkSttlmDt: (*common.ISODate)(&message.FedNowMsg.OriginalIdentifier.CreationDateTime),
					RtrdIntrBkSttlmAmt: pacs004.ActiveCurrencyAndAmount{
						Ccy:  pacs004.ActiveCurrencyCode(message.FedNowMsg.PaymentReturn.ReturnedAmount.Ccy),
						Text: string(message.FedNowMsg.PaymentReturn.ReturnedAmount.Text),
					},
					IntrBkSttlmDt: (*common.ISODate)(&message.FedNowMsg.CreationDateTime),
					ChrgBr:        &chargebearer,
					InstgAgt: &pacs004.BranchAndFinancialInstitutionIdentification6{
						FinInstnId: pacs004.FinancialInstitutionIdentification18{
							ClrSysMmbId: &pacs004.ClearingSystemMemberIdentification2{
								ClrSysId: &pacs004.ClearingSystemIdentification2Choice{
									Cd: &clearingSystemId,
								},
								MmbId: pacs004.Max35Text(message.FedNowMsg.SenderDI.SenderABANumber),
							},
						},
					},
					InstdAgt: &pacs004.BranchAndFinancialInstitutionIdentification6{
						FinInstnId: pacs004.FinancialInstitutionIdentification18{
							ClrSysMmbId: &pacs004.ClearingSystemMemberIdentification2{
								ClrSysId: &pacs004.ClearingSystemIdentification2Choice{
									Cd: &clearingSystemId,
								},
								MmbId: pacs004.Max35Text(message.FedNowMsg.ReceiverDI.ReceiverABANumber),
							},
						},
					},
					RtrChain: &pacs004.TransactionParties8{
						Dbtr: pacs004.Party40Choice{
							Pty: &pacs004.PartyIdentification135{
								Nm: (*pacs004.Max140Text)(message.FedNowMsg.Originator.Personal.Name),
								PstlAdr: &pacs004.PostalAddress24{
									StrtNm:      (*pacs004.Max70Text)(message.FedNowMsg.Originator.Personal.Address.StreetName),
									BldgNb:      (*pacs004.Max16Text)(message.FedNowMsg.Originator.Personal.Address.BuildingNumber),
									TwnNm:       (*pacs004.Max35Text)(message.FedNowMsg.Originator.Personal.Address.TownName),
									CtrySubDvsn: (*pacs004.Max35Text)(message.FedNowMsg.Originator.Personal.Address.CountrySubdivision),
									PstCd:       (*pacs004.Max16Text)(message.FedNowMsg.Originator.Personal.Address.PostalCode),
									Ctry:        (*pacs004.CountryCode)(message.FedNowMsg.Originator.Personal.Address.Country),
								},
							},
						},
						DbtrAcct: &pacs004.CashAccount38{
							Id: pacs004.AccountIdentification4Choice{
								Othr: &pacs004.GenericAccountIdentification1{
									Id: pacs004.Max34Text(message.FedNowMsg.Originator.Personal.Identifier),
								},
							},
						},
						DbtrAgt: &pacs004.BranchAndFinancialInstitutionIdentification6{
							FinInstnId: pacs004.FinancialInstitutionIdentification18{
								ClrSysMmbId: &pacs004.ClearingSystemMemberIdentification2{
									MmbId: pacs004.Max35Text(message.FedNowMsg.SenderDI.SenderABANumber),
									ClrSysId: &pacs004.ClearingSystemIdentification2Choice{
										Cd: &clearingSystemId,
									},
								},
								Nm: (*pacs004.Max140Text)(message.FedNowMsg.SenderDI.Name),
							},
						},
						CdtrAgt: &pacs004.BranchAndFinancialInstitutionIdentification6{
							FinInstnId: pacs004.FinancialInstitutionIdentification18{
								ClrSysMmbId: &pacs004.ClearingSystemMemberIdentification2{
									MmbId: pacs004.Max35Text(message.FedNowMsg.ReceiverDI.ReceiverABANumber),
									ClrSysId: &pacs004.ClearingSystemIdentification2Choice{
										Cd: &clearingSystemId,
									},
								},
								Nm: (*pacs004.Max140Text)(message.FedNowMsg.ReceiverDI.Name),
							},
						},
						Cdtr: pacs004.Party40Choice{
							Pty: &pacs004.PartyIdentification135{
								Nm: (*pacs004.Max140Text)(message.FedNowMsg.Beneficiary.Personal.Name),
								PstlAdr: &pacs004.PostalAddress24{
									StrtNm:      (*pacs004.Max70Text)(message.FedNowMsg.Beneficiary.Personal.Address.StreetName),
									BldgNb:      (*pacs004.Max16Text)(message.FedNowMsg.Beneficiary.Personal.Address.BuildingNumber),
									TwnNm:       (*pacs004.Max35Text)(message.FedNowMsg.Beneficiary.Personal.Address.TownName),
									CtrySubDvsn: (*pacs004.Max35Text)(message.FedNowMsg.Beneficiary.Personal.Address.CountrySubdivision),
									PstCd:       (*pacs004.Max16Text)(message.FedNowMsg.Beneficiary.Personal.Address.PostalCode),
									Ctry:        (*pacs004.CountryCode)(message.FedNowMsg.Beneficiary.Personal.Address.Country),
								},
							},
						},
						CdtrAcct: &pacs004.CashAccount38{
							Id: pacs004.AccountIdentification4Choice{
								Othr: &pacs004.GenericAccountIdentification1{
									Id: pacs004.Max34Text(message.FedNowMsg.Beneficiary.Personal.Identifier),
								},
							},
						},
					},
					RtrRsnInf: []pacs004.PaymentReturnReason6{
						{
							Rsn: &pacs004.ReturnReason5Choice{
								Cd: message.FedNowMsg.PaymentReturn.ReturnReason,
							},
							AddtlInf: []pacs004.Max105Text{
								pacs004.Max105Text(*message.FedNowMsg.PaymentReturn.AdditionalInformation),
							},
						},
					},
					OrgnlTxRef: &pacs004.OriginalTransactionReference32{
						PmtTpInf: &pacs004.PaymentTypeInformation27{
							LclInstrm: &pacs004.LocalInstrument2Choice{
								Prtry: &localInstrument,
							},
						},
					},
				},
			},
		},
	}

	return pacsDoc, nil
}
