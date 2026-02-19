package admi

import (
	admi_011_001_01 "github.com/mbanq/iso20022-go/ISO20022/admi_011_001_01"
	head "github.com/mbanq/iso20022-go/ISO20022/head_001_001_02"
	"github.com/mbanq/iso20022-go/pkg/common"
)

// ParseAdmi011 parses an incoming admi.011 SystemEventAcknowledgement message.
// This is FedNow's response to our admi.004 ParticipantBroadcast messages (PING, FPON, FPOF, etc.).
func ParseAdmi011(appHdr head.BusinessApplicationHeaderV02, document admi_011_001_01.Document) (*FedNowMessageSystemResponse, error) {
	sysEvtAck := document.SysEvtAck

	// Extract original reference (the admi.004 message ID we sent)
	var originalReference string
	if sysEvtAck.OrgtrRef != nil {
		originalReference = string(*sysEvtAck.OrgtrRef)
	}

	// Extract acknowledgement details
	var eventCode string
	var eventParam string
	var eventTime *common.ISODateTime

	if sysEvtAck.AckDtls != nil {
		eventCode = string(sysEvtAck.AckDtls.EvtCd)

		// Extract first event parameter (connection identifier)
		if len(sysEvtAck.AckDtls.EvtParam) > 0 {
			eventParam = string(sysEvtAck.AckDtls.EvtParam[0])
		}

		eventTime = sysEvtAck.AckDtls.EvtTm
	}

	fednowMsg := FedNowMessageSystemResponse{
		FedNowMsg: FedNowSystemResponse{
			CreationDateTime:  common.ISODateTime(appHdr.CreDt),
			MessageID:         string(sysEvtAck.MsgId),
			OriginalReference: originalReference,
			EventCode:         eventCode,
			EventParam:        eventParam,
			EventTime:         eventTime,
		},
	}

	return &fednowMsg, nil
}
