package pacs

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mbanq/iso20022-go/ISO20022/pacs_002_001_10"
	"github.com/mbanq/iso20022-go/ISO20022/pacs_004_001_10"
	"github.com/mbanq/iso20022-go/ISO20022/pacs_008_001_08"
	"github.com/mbanq/iso20022-go/pkg/common"
)

type FedNowMessageCCT struct {
	FedNowMsg FedNowDetails `json:"fedNowMessage"`
}

func (f FedNowMessageCCT) IsFedNowMessage() {}

type FedNowMessageACK struct {
	FedNowMsg FedNowACK `json:"fedNowMessage"`
}

func (f FedNowMessageACK) IsFedNowMessage() {}

type FedNowMessageRtn struct {
	FedNowMsg FedNowRtn `json:"fedNowMessage"`
}

func (f FedNowMessageRtn) IsFedNowMessage() {}

type FedNowDetails struct {
	CreationDateTime common.ISODateTime          `json:"creationDateTime"`
	Identifier       FedNowIdentifier            `json:"identifier"`
	PaymentType      FedNowPaymentType           `json:"paymentType"`
	Amount           FedNowAmount                `json:"amount"`
	SenderDI         FedNowDepositoryInstitution `json:"senderDepositoryInstitution"`
	ReceiverDI       FedNowDepositoryInstitution `json:"receiverDepositoryInstitution"`
	Originator       FedNowParty                 `json:"originator"`
	Beneficiary      FedNowParty                 `json:"beneficiary"`
}

type FedNowACK struct {
	CreationDateTime   common.ISODateTime          `json:"creationDateTime"`
	Identifier         FedNowIdentifier            `json:"identifier"`
	OriginalIdentifier FedNowIdentifier            `json:"originalIdentifier"`
	PaymentStatus      PaymentStatus               `json:"paymentStatus"`
	SenderDI           FedNowDepositoryInstitution `json:"senderDepositoryInstitution"`
	ReceiverDI         FedNowDepositoryInstitution `json:"receiverDepositoryInstitution"`
}

type FedNowRtn struct {
	CreationDateTime   common.ISODateTime          `json:"creationDateTime"`
	Identifier         FedNowIdentifier            `json:"identifier"`
	OriginalIdentifier FedNowIdentifier            `json:"originalIdentifier"`
	Amount             FedNowAmount                `json:"amount"`
	PaymentReturn      PaymentReturn               `json:"paymentReturn"`
	SenderDI           FedNowDepositoryInstitution `json:"senderDepositoryInstitution"`
	ReceiverDI         FedNowDepositoryInstitution `json:"receiverDepositoryInstitution"`
	Originator         FedNowParty                 `json:"originator"`
	Beneficiary        FedNowParty                 `json:"beneficiary"`
}

type FedNowIdentifier struct {
	BusinessMessageID pacs_008_001_08.Max35Text         `json:"businessMessageId"`
	MessageID         pacs_008_001_08.Max35Text         `json:"messageId"`
	MessageType       pacs_008_001_08.Max35Text         `json:"messageType,omitempty"`
	InstructionID     *pacs_008_001_08.Max35Text        `json:"instructionId"`
	EndToEndID        pacs_008_001_08.Max35Text         `json:"endToEndId,omitempty"`
	TransactionID     *pacs_008_001_08.Max35Text        `json:"transactionId,omitempty"`
	UETR              *pacs_008_001_08.UUIDv4Identifier `json:"uetr,omitempty"`
	CreationDateTime  common.ISODateTime                `json:"creationDateTime,omitempty"`
}

type FedNowPaymentType struct {
	CategoryPurpose *pacs_008_001_08.ExternalCategoryPurpose1Code `json:"categoryPurpose"`
}

type FedNowAmount struct {
	Text json.Number                        `json:"amount"`
	Ccy  pacs_008_001_08.ActiveCurrencyCode `json:"currency"`
}

type FedNowDepositoryInstitution struct {
	SenderABANumber   pacs_008_001_08.Max35Text   `json:"senderABANumber,omitempty"`
	ReceiverABANumber pacs_008_001_08.Max35Text   `json:"receiverABANumber,omitempty"`
	Name              *pacs_008_001_08.Max140Text `json:"senderShortName"`
}

type FedNowParty struct {
	Personal FedNowPersonal `json:"personal"`
}

type FedNowPersonal struct {
	Name       *pacs_008_001_08.Max140Text `json:"name"`
	Address    FedNowPstlAdr               `json:"postalAddress"`
	Identifier pacs_008_001_08.Max34Text   `json:"identifier"`
}

type FedNowPstlAdr struct {
	StreetName         *pacs_008_001_08.Max70Text   `json:"StreetName"`
	BuildingNumber     *pacs_008_001_08.Max16Text   `json:"BuildingNumber"`
	PostBox            *pacs_008_001_08.Max16Text   `json:"PostBox"`
	TownName           *pacs_008_001_08.Max35Text   `json:"TownName"`
	CountrySubdivision *pacs_008_001_08.Max35Text   `json:"CountrySubDivision"`
	PostalCode         *pacs_008_001_08.Max16Text   `json:"PostalCode"`
	Country            *pacs_008_001_08.CountryCode `json:"Country"`
}

type PaymentStatus struct {
	//TODO: Add Optional Field - Originator
	PaymentStatus         *pacs_002_001_10.ExternalPaymentTransactionStatus1Code `json:"paymentStatus"`
	AcceptanceDateTime    *common.ISODateTime                                    `json:"acceptanceDateTime,omitempty"`
	StatusReason          *pacs_002_001_10.ExternalStatusReason1Code             `json:"statusReason,omitempty"`
	AdditionalInformation *pacs_002_001_10.Max105Text                            `json:"additionalInformation,omitempty"`
}

type PaymentReturn struct {
	ReturnReason          *pacs_004_001_10.ExternalReturnReason1Code `json:"returnReason"`
	AdditionalInformation *pacs_004_001_10.Max105Text                `json:"additionalInformation,omitempty"`
	ReturnedAmount        FedNowAmount                               `json:"returnedAmount"`
}

func (address FedNowPstlAdr) ValidateAddress() error {
	var missingFields []string
	if address.StreetName == nil || *address.StreetName == "" {
		missingFields = append(missingFields, "StreetName")
	}
	if address.TownName == nil || *address.TownName == "" {
		missingFields = append(missingFields, "TownName")
	}
	if address.CountrySubdivision == nil || *address.CountrySubdivision == "" {
		missingFields = append(missingFields, "CountrySubdivision")
	}
	if address.PostalCode == nil || *address.PostalCode == "" {
		missingFields = append(missingFields, "PostalCode")
	}
	if address.Country == nil || *address.Country == "" {
		missingFields = append(missingFields, "Country")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("missing required address fields: %s", strings.Join(missingFields, ", "))
	}

	return nil
}
