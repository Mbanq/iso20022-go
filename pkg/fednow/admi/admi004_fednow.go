package admi

import (
	admi_004_001_02 "github.com/mbanq/iso20022-go/ISO20022/admi_004_001_02"
	head "github.com/mbanq/iso20022-go/ISO20022/head_001_001_02"
	"github.com/mbanq/iso20022-go/pkg/common"
)

// ParseFedNowBroadcast parses an incoming admi.004 FedNowBroadcast message.
// This is sent by FedNow to all participants to notify of system events or other participants' status changes.
// Events include: FPON, FPOF (other participants), ROLL, EXTN, FNKY, etc.
func ParseFedNowBroadcast(appHdr head.BusinessApplicationHeaderV02, document admi_004_001_02.Document) (*FedNowMessageFedNowBroadcast, error) {
	sysEvtNtfctn := document.SysEvtNtfctn
	evtInf := sysEvtNtfctn.EvtInf

	// Extract event parameters (can be multiple for FedNow broadcasts)
	var eventParams []string
	for _, param := range evtInf.EvtParam {
		eventParams = append(eventParams, string(param))
	}

	// Extract event description
	var eventDescription *string
	if evtInf.EvtDesc != nil {
		desc := string(*evtInf.EvtDesc)
		eventDescription = &desc
	}

	fednowMsg := FedNowMessageFedNowBroadcast{
		FedNowMsg: FedNowFedNowBroadcast{
			CreationDateTime: common.ISODateTime(appHdr.CreDt),
			MessageID:        string(appHdr.BizMsgIdr),
			EventCode:        string(evtInf.EvtCd),
			EventParams:      eventParams,
			EventDescription: eventDescription,
			EventTime:        evtInf.EvtTm,
		},
	}

	return &fednowMsg, nil
}
