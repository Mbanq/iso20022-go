package pain

import (
	"encoding/json"
	"encoding/xml"
	"fmt"

	"github.com/mbanq/iso20022-go/ISO20022/head_001_001_02"
	"github.com/mbanq/iso20022-go/ISO20022/pain_013_001_07"
	"github.com/mbanq/iso20022-go/pkg/common"
	"github.com/mbanq/iso20022-go/pkg/fednow/config"
)

func BuildPain013Struct(message FedNowMessageRFP, msgConfig *config.Config) (*pain_013_001_07.Document, error) {

	fedMsg := message.FedNowMsg

	// Assigning Configuration Values
	localInstrument := pain_013_001_07.Max35Text(*msgConfig.LocalInstrument.Prtry)
	clearingSystemId := pain_013_001_07.ExternalClearingSystemIdentification1Code(msgConfig.ClearingSystemId)

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

	// Building the Pain013 Struct
	painDoc := &pain_013_001_07.Document{
		XMLName: xml.Name{Space: "urn:iso:std:iso:20022:tech:xsd:pain.013.001.07", Local: "Document"},
		CdtrPmtActvtnReq: pain_013_001_07.CreditorPaymentActivationRequestV07{
			GrpHdr: pain_013_001_07.GroupHeader78{
				MsgId:   fedMsg.Identifier.MessageID,
				CreDtTm: fedMsg.ExecutionInfo.ExecutionDate,
				NbOfTxs: pain_013_001_07.Max15NumericText("1"),
				InitgPty: pain_013_001_07.PartyIdentification135{
					Nm: fedMsg.ExecutionInfo.InitiatingParty,
				},
			},
			PmtInf: []pain_013_001_07.PaymentInstruction31{
				{
					PmtInfId:    &fedMsg.Identifier.TransactionID,
					PmtMtd:      "TRF",
					ReqdExctnDt: pain_013_001_07.DateAndDateTime2Choice{DtTm: &fedMsg.ExecutionInfo.ExecutionDate},
					XpryDt:      &pain_013_001_07.DateAndDateTime2Choice{DtTm: &fedMsg.ExecutionInfo.ExpiryDate},
					Dbtr: pain_013_001_07.PartyIdentification135{
						Nm: fedMsg.Beneficiary.Personal.Name,
						PstlAdr: &pain_013_001_07.PostalAddress24{
							StrtNm:      fedMsg.Beneficiary.Personal.Address.StreetName,
							BldgNb:      fedMsg.Beneficiary.Personal.Address.BuildingNumber,
							PstCd:       fedMsg.Beneficiary.Personal.Address.PostalCode,
							TwnNm:       fedMsg.Beneficiary.Personal.Address.TownName,
							CtrySubDvsn: fedMsg.Beneficiary.Personal.Address.CountrySubdivision,
							Ctry:        fedMsg.Beneficiary.Personal.Address.Country,
						},
					},
					DbtrAcct: &pain_013_001_07.CashAccount38{
						Id: pain_013_001_07.AccountIdentification4Choice{
							Othr: &pain_013_001_07.GenericAccountIdentification1{
								Id: fedMsg.Beneficiary.Personal.Identifier,
							},
						},
					},
					DbtrAgt: pain_013_001_07.BranchAndFinancialInstitutionIdentification6{
						FinInstnId: pain_013_001_07.FinancialInstitutionIdentification18{
							ClrSysMmbId: &pain_013_001_07.ClearingSystemMemberIdentification2{
								MmbId: fedMsg.ReceiverDI.ReceiverABANumber,
								ClrSysId: &pain_013_001_07.ClearingSystemIdentification2Choice{
									Cd: &clearingSystemId,
								},
							},
						},
					},
					CdtTrfTx: []pain_013_001_07.CreditTransferTransaction35{
						{
							PmtId: pain_013_001_07.PaymentIdentification6{
								EndToEndId: fedMsg.Identifier.TransactionID,
							},
							PmtTpInf: &pain_013_001_07.PaymentTypeInformation26{
								LclInstrm: &pain_013_001_07.LocalInstrument2Choice{
									Prtry: &localInstrument,
								},
								CtgyPurp: &pain_013_001_07.CategoryPurpose1Choice{
									Prtry: &fedMsg.PaymentType.CategoryPurpose,
								},
							},
							Amt: pain_013_001_07.AmountType4Choice{
								InstdAmt: &pain_013_001_07.ActiveOrHistoricCurrencyAndAmount{
									Ccy:  fedMsg.Amount.Ccy,
									Text: fmt.Sprintf("%.2f", amountFloat),
								},
							},
							ChrgBr: pain_013_001_07.ChargeBearerType1Code(msgConfig.ChargeBearer),
							CdtrAgt: pain_013_001_07.BranchAndFinancialInstitutionIdentification6{
								FinInstnId: pain_013_001_07.FinancialInstitutionIdentification18{
									ClrSysMmbId: &pain_013_001_07.ClearingSystemMemberIdentification2{
										MmbId: fedMsg.SenderDI.SenderABANumber,
										ClrSysId: &pain_013_001_07.ClearingSystemIdentification2Choice{
											Cd: &clearingSystemId,
										},
									},
								},
							},
							Cdtr: pain_013_001_07.PartyIdentification135{
								Nm: fedMsg.Originator.Personal.Name,
								PstlAdr: &pain_013_001_07.PostalAddress24{
									StrtNm:      fedMsg.Originator.Personal.Address.StreetName,
									BldgNb:      fedMsg.Originator.Personal.Address.BuildingNumber,
									PstCd:       fedMsg.Originator.Personal.Address.PostalCode,
									TwnNm:       fedMsg.Originator.Personal.Address.TownName,
									CtrySubDvsn: fedMsg.Originator.Personal.Address.CountrySubdivision,
									Ctry:        fedMsg.Originator.Personal.Address.Country,
								},
							},
							CdtrAcct: &pain_013_001_07.CashAccount38{
								Id: pain_013_001_07.AccountIdentification4Choice{
									Othr: &pain_013_001_07.GenericAccountIdentification1{
										Id: fedMsg.Originator.Personal.Identifier,
									},
								},
							},
						}},
				},
			},
		},
	}
	return painDoc, nil
}

func ParsePain013(appHdr head_001_001_02.BusinessApplicationHeaderV02, document pain_013_001_07.Document) (*FedNowMessageRFP, error) {

	payment_request := document.CdtrPmtActvtnReq

	fednowMsg := FedNowMessageRFP{
		FedNowMsg: FedNowDetails{
			CreationDateTime: common.ISODateTime(appHdr.CreDt),
			Identifier: FedNowIdentifier{
				BusinessMessageID: pain_013_001_07.Max35Text(appHdr.BizMsgIdr),
				MessageID:         pain_013_001_07.Max35Text(payment_request.GrpHdr.MsgId),
				InstructionID:     *payment_request.PmtInf[0].PmtInfId,
				EndToEndID:        payment_request.PmtInf[0].CdtTrfTx[0].PmtId.EndToEndId,
				TransactionID:     *payment_request.PmtInf[0].PmtInfId,
			},
			PaymentType: FedNowPaymentType{
				CategoryPurpose: *payment_request.PmtInf[0].CdtTrfTx[0].PmtTpInf.CtgyPurp.Prtry,
			},
			ExecutionInfo: FedNowExecutionInfo{
				InitiatingParty: payment_request.GrpHdr.InitgPty.Nm,
			},
			Amount: FedNowAmount{
				Text: json.Number(payment_request.PmtInf[0].CdtTrfTx[0].Amt.InstdAmt.Text),
				Ccy:  payment_request.PmtInf[0].CdtTrfTx[0].Amt.InstdAmt.Ccy,
			},
			SenderDI: FedNowDepositoryInstitution{
				ReceiverABANumber: pain_013_001_07.Max35Text(appHdr.To.FIId.FinInstnId.ClrSysMmbId.MmbId),
			},
			ReceiverDI: FedNowDepositoryInstitution{
				SenderABANumber: pain_013_001_07.Max35Text(appHdr.Fr.FIId.FinInstnId.ClrSysMmbId.MmbId),
			},
			Originator: FedNowParty{
				Personal: FedNowPersonal{
					Name: payment_request.PmtInf[0].Dbtr.Nm,
					Address: FedNowPstlAdr{
						StreetName:         payment_request.PmtInf[0].Dbtr.PstlAdr.StrtNm,
						BuildingNumber:     payment_request.PmtInf[0].Dbtr.PstlAdr.BldgNb,
						TownName:           payment_request.PmtInf[0].Dbtr.PstlAdr.TwnNm,
						CountrySubdivision: payment_request.PmtInf[0].Dbtr.PstlAdr.CtrySubDvsn,
						PostalCode:         payment_request.PmtInf[0].Dbtr.PstlAdr.PstCd,
						Country:            payment_request.PmtInf[0].Dbtr.PstlAdr.Ctry,
					},
					Identifier: pain_013_001_07.Max34Text(payment_request.PmtInf[0].DbtrAcct.Id.Othr.Id),
				},
			},
			Beneficiary: FedNowParty{
				Personal: FedNowPersonal{
					Name: payment_request.PmtInf[0].CdtTrfTx[0].Cdtr.Nm,
					Address: FedNowPstlAdr{
						StreetName:         payment_request.PmtInf[0].CdtTrfTx[0].Cdtr.PstlAdr.StrtNm,
						BuildingNumber:     payment_request.PmtInf[0].CdtTrfTx[0].Cdtr.PstlAdr.BldgNb,
						TownName:           payment_request.PmtInf[0].CdtTrfTx[0].Cdtr.PstlAdr.TwnNm,
						CountrySubdivision: payment_request.PmtInf[0].CdtTrfTx[0].Cdtr.PstlAdr.CtrySubDvsn,
						PostalCode:         payment_request.PmtInf[0].CdtTrfTx[0].Cdtr.PstlAdr.PstCd,
						Country:            payment_request.PmtInf[0].CdtTrfTx[0].Cdtr.PstlAdr.Ctry,
					},
					Identifier: pain_013_001_07.Max34Text(payment_request.PmtInf[0].CdtTrfTx[0].CdtrAcct.Id.Othr.Id),
				},
			},
		},
	}

	if payment_request.PmtInf[0].ReqdExctnDt.DtTm != nil {
		fednowMsg.FedNowMsg.ExecutionInfo.ExecutionDate = common.ISODateTime(*payment_request.PmtInf[0].ReqdExctnDt.DtTm)
	}
	if payment_request.PmtInf[0].XpryDt.DtTm != nil {
		fednowMsg.FedNowMsg.ExecutionInfo.ExpiryDate = common.ISODateTime(*payment_request.PmtInf[0].XpryDt.DtTm)
	}

	if payment_request.GrpHdr.InitgPty.PstlAdr != nil {
		fednowMsg.FedNowMsg.ExecutionInfo.InitiatingPartyAddress = FedNowPstlAdr{
			StreetName:         payment_request.GrpHdr.InitgPty.PstlAdr.StrtNm,
			BuildingNumber:     payment_request.GrpHdr.InitgPty.PstlAdr.BldgNb,
			TownName:           payment_request.GrpHdr.InitgPty.PstlAdr.TwnNm,
			CountrySubdivision: payment_request.GrpHdr.InitgPty.PstlAdr.CtrySubDvsn,
			PostalCode:         payment_request.GrpHdr.InitgPty.PstlAdr.PstCd,
			Country:            payment_request.GrpHdr.InitgPty.PstlAdr.Ctry,
		}
	}

	return &fednowMsg, nil
}
