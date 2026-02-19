package admi

import (
	"encoding/json"
	"encoding/xml"
	"time"

	admi_004_001_02 "github.com/mbanq/iso20022-go/ISO20022/admi_004_001_02"
	"github.com/mbanq/iso20022-go/pkg/common"
	"github.com/mbanq/iso20022-go/pkg/fednow/config"
)

// BuildParticipantBroadcast creates an outgoing admi.004 ParticipantBroadcast message.
// Used for PING, FPON (sign on), FPOF (sign off), FPCR (reconnect), FPCD (disconnect).
func BuildParticipantBroadcast(message FedNowMessageParticipantBroadcast, msgConfig *config.Config) (*admi_004_001_02.Document, error) {
	fedMsg := message.FedNowMsg

	// Set event time to now if not provided
	eventTime := fedMsg.EventTime
	if eventTime == nil || time.Time(*eventTime).IsZero() {
		now := common.ISODateTime(time.Now().UTC())
		eventTime = &now
	}

	// Build event parameters (single connection identifier for participant broadcast)
	var evtParam []admi_004_001_02.Max35Text
	if fedMsg.ConnectionIdentifier != "" {
		evtParam = []admi_004_001_02.Max35Text{admi_004_001_02.Max35Text(fedMsg.ConnectionIdentifier)}
	}

	admiDoc := &admi_004_001_02.Document{
		XMLName: xml.Name{Space: "urn:iso:std:iso:20022:tech:xsd:admi.004.001.02", Local: "Document"},
		SysEvtNtfctn: admi_004_001_02.SystemEventNotificationV02{
			EvtInf: admi_004_001_02.Event2{
				EvtCd:    admi_004_001_02.Max4AlphaNumericText(fedMsg.EventCode),
				EvtParam: evtParam,
				EvtDesc:  convertToMax1000Text(fedMsg.EventDescription),
				EvtTm:    eventTime,
			},
		},
	}

	return admiDoc, nil
}

// BuildParticipantBroadcastFromJSON creates an admi.004 Document from JSON payload.
func BuildParticipantBroadcastFromJSON(payload []byte, cfg *config.Config) (*admi_004_001_02.Document, error) {
	var message FedNowMessageParticipantBroadcast
	if err := json.Unmarshal(payload, &message); err != nil {
		return nil, err
	}
	return BuildParticipantBroadcast(message, cfg)
}

// NewFPONMessage creates a sign-on (Function Parameter ON) message.
// connectionIdentifier can be either a 9-digit RTN or 7-character connection point ID.
func NewFPONMessage(messageID string, connectionIdentifier string) FedNowMessageParticipantBroadcast {
	now := common.ISODateTime(time.Now().UTC())
	return FedNowMessageParticipantBroadcast{
		FedNowMsg: FedNowParticipantBroadcast{
			MessageID:            messageID,
			EventCode:            "FPON",
			ConnectionIdentifier: connectionIdentifier,
			EventDescription:     nil,
			EventTime:            &now,
		},
	}
}

// NewFPOFMessage creates a sign-off (Function Parameter OFF) message.
// connectionIdentifier can be either a 9-digit RTN or 7-character connection point ID.
func NewFPOFMessage(messageID string, connectionIdentifier string) FedNowMessageParticipantBroadcast {
	now := common.ISODateTime(time.Now().UTC())
	return FedNowMessageParticipantBroadcast{
		FedNowMsg: FedNowParticipantBroadcast{
			MessageID:            messageID,
			EventCode:            "FPOF",
			ConnectionIdentifier: connectionIdentifier,
			EventDescription:     nil,
			EventTime:            &now,
		},
	}
}

// NewPINGMessage creates a connectivity check (PING) message.
func NewPINGMessage(messageID string) FedNowMessageParticipantBroadcast {
	now := common.ISODateTime(time.Now().UTC())
	return FedNowMessageParticipantBroadcast{
		FedNowMsg: FedNowParticipantBroadcast{
			MessageID:            messageID,
			EventCode:            "PING",
			ConnectionIdentifier: "", // PING doesn't require connection identifier
			EventDescription:     nil,
			EventTime:            &now,
		},
	}
}

// Helper to convert string pointer to Max1000Text pointer
func convertToMax1000Text(s *string) *admi_004_001_02.Max1000Text {
	if s == nil {
		return nil
	}
	text := admi_004_001_02.Max1000Text(*s)
	return &text
}
