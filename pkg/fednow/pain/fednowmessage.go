package pain

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mbanq/iso20022-go/ISO20022/pain_013_001_07"
	"github.com/mbanq/iso20022-go/pkg/common"
)

type FedNowMessageRFP struct {
	FedNowMsg FedNowDetails `json:"fedNowMessage"`
}

func (f FedNowMessageRFP) IsFedNowMessage() {}

type FedNowDetails struct {
	CreationDateTime common.ISODateTime          `json:"creationDateTime"`
	Identifier       FedNowIdentifier            `json:"identifier"`
	PaymentType      FedNowPaymentType           `json:"paymentType"`
	ExecutionInfo    FedNowExecutionInfo         `json:"executionInfo"`
	Amount           FedNowAmount                `json:"amount"`
	SenderDI         FedNowDepositoryInstitution `json:"senderDepositoryInstitution"`
	ReceiverDI       FedNowDepositoryInstitution `json:"receiverDepositoryInstitution"`
	Originator       FedNowParty                 `json:"originator"`
	Beneficiary      FedNowParty                 `json:"beneficiary"`
}

type FedNowIdentifier struct {
	BusinessMessageID pain_013_001_07.Max35Text `json:"businessMessageId"`
	MessageID         pain_013_001_07.Max35Text `json:"messageId"`
	MessageType       pain_013_001_07.Max35Text `json:"messageType"`
	InstructionID     pain_013_001_07.Max35Text `json:"instructionId"`
	TransactionID     pain_013_001_07.Max35Text `json:"transactionId"`
	UETR              pain_013_001_07.Max35Text `json:"uetr"`
	EndToEndID        pain_013_001_07.Max35Text `json:"endToEndId"`
}

type FedNowPaymentType struct {
	CategoryPurpose pain_013_001_07.Max35Text `json:"categoryPurpose"`
}

type FedNowExecutionInfo struct {
	InitiatingParty        *pain_013_001_07.Max140Text `json:"initiatingParty"`
	InitiatingPartyAddress FedNowPstlAdr               `json:"initiatingPartyAddress"`
	ExecutionDate          common.ISODateTime          `json:"executionDate"`
	ExpiryDate             common.ISODateTime          `json:"expiryDate"`
}

type FedNowAmount struct {
	Text json.Number                                  `json:"amount"`
	Ccy  pain_013_001_07.ActiveOrHistoricCurrencyCode `json:"currency"`
}

type FedNowDepositoryInstitution struct {
	SenderABANumber   pain_013_001_07.Max35Text   `json:"senderABANumber,omitempty"`
	ReceiverABANumber pain_013_001_07.Max35Text   `json:"receiverABANumber,omitempty"`
	Name              *pain_013_001_07.Max140Text `json:"senderShortName"`
}

type FedNowParty struct {
	Personal FedNowPersonal `json:"personal"`
}

type FedNowPersonal struct {
	Name       *pain_013_001_07.Max140Text `json:"name"`
	Address    FedNowPstlAdr               `json:"postalAddress"`
	Identifier pain_013_001_07.Max34Text   `json:"identifier"`
}

type FedNowPstlAdr struct {
	StreetName         *pain_013_001_07.Max70Text   `json:"StreetName"`
	BuildingNumber     *pain_013_001_07.Max16Text   `json:"BuildingNumber"`
	PostBox            *pain_013_001_07.Max16Text   `json:"PostalBox"`
	TownName           *pain_013_001_07.Max35Text   `json:"TownName"`
	CountrySubdivision *pain_013_001_07.Max35Text   `json:"CountrySubDivision"`
	PostalCode         *pain_013_001_07.Max16Text   `json:"PostalCode"`
	Country            *pain_013_001_07.CountryCode `json:"Country"`
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
