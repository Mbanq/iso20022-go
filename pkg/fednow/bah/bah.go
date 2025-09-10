package bah

import (
	"time"

	bah "github.com/Mbanq/iso20022-go/ISO20022/head_001_001_02"
	"github.com/Mbanq/iso20022-go/pkg/common"
	"github.com/Mbanq/iso20022-go/pkg/fednow/config"
)

func BuildBah(messageId string, msgConfig *config.Config, msgType string) (*bah.BusinessApplicationHeaderV02, error) {

	now := time.Now().In(common.EstLocation)

	bahMsg := &bah.BusinessApplicationHeaderV02{
		Fr: bah.Party44Choice{
			FIId: &bah.BranchAndFinancialInstitutionIdentification6{
				FinInstnId: bah.FinancialInstitutionIdentification18{
					ClrSysMmbId: &bah.ClearingSystemMemberIdentification2{
						MmbId: bah.Max35Text(msgConfig.IspId),
					},
				},
			},
		},
		To: bah.Party44Choice{
			FIId: &bah.BranchAndFinancialInstitutionIdentification6{
				FinInstnId: bah.FinancialInstitutionIdentification18{
					ClrSysMmbId: &bah.ClearingSystemMemberIdentification2{
						MmbId: bah.Max35Text(msgConfig.FrbId),
					},
				},
			},
		},
		BizMsgIdr: bah.Max35Text(messageId),
		MsgDefIdr: bah.Max35Text(msgType),
		MktPrctc: &bah.ImplementationSpecification1{
			Regy: bah.Max350Text("www2.swift.com/mystandards/#/group/Federal_Reserve_Financial_Services/FedNow_Service"),
			Id:   bah.Max2048Text("frb.fednow.01"),
		},
		CreDt: (common.ISODateTime)(now),
	}

	return bahMsg, nil
}
