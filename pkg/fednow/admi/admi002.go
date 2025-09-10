package admi

import (
	admi "github.com/Mbanq/iso20022-go/ISO20022/admi_002_001_01"
	head "github.com/Mbanq/iso20022-go/ISO20022/head_001_001_02"
	"github.com/Mbanq/iso20022-go/pkg/common"
)

func ParseAdmi002Struct(admiDoc *admi.Document, appHdr head.BusinessApplicationHeaderV02) (FedNowMessageADM, error) {
	fedMsg := FedNowMessageADM{
		FedNowMsg: FedNowADM{
			CreationDateTime: common.ISODateTime(appHdr.CreDt),
			Identifier: FedNowIdentifier{
				BusinessMessageID: appHdr.BizMsgIdr,
				MessageType:       appHdr.MsgDefIdr,
				MessageID:         appHdr.BizMsgIdr,
			},
			Reference: admiDoc.Admi00200101.RltdRef.Ref,
			Reason: RejectionReason{
				RejectionReason:   admiDoc.Admi00200101.Rsn.RjctgPtyRsn,
				RejectionDateTime: admiDoc.Admi00200101.Rsn.RjctnDtTm,
			},
		},
	}
	return fedMsg, nil
}
