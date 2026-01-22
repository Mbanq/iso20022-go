package admi

import (
	"encoding/json"
	"encoding/xml"
	"time"

	admi_007_001_01 "github.com/mbanq/iso20022-go/ISO20022/admi_007_001_01"
	head "github.com/mbanq/iso20022-go/ISO20022/head_001_001_02"
	"github.com/mbanq/iso20022-go/pkg/common"
	"github.com/mbanq/iso20022-go/pkg/fednow/config"
)

func BuildAdmi007Struct(message FedNowMessageRctAck, msgConfig *config.Config) (*admi_007_001_01.Document, error) {
	fedMsg := message.FedNowMsg

	var creationTimePtr *common.ISODateTime
	if !time.Time(fedMsg.CreationDateTime).IsZero() {
		value := fedMsg.CreationDateTime
		creationTimePtr = &value
	}

	msgHeader := admi_007_001_01.MessageHeader10{
		MsgId:   admi_007_001_01.Max35Text(fedMsg.Identifier.MessageID),
		CreDtTm: creationTimePtr,
		QryNm:   fedMsg.QueryName,
	}

	reports := make([]admi_007_001_01.ReceiptAcknowledgementReport2, 0, len(fedMsg.Reports))
	for _, report := range fedMsg.Reports {
		reports = append(reports, admi_007_001_01.ReceiptAcknowledgementReport2{
			RltdRef: admi_007_001_01.MessageReference1{
				Ref:     report.RelatedReference.Reference,
				MsgNm:   report.RelatedReference.MessageName,
				RefIssr: report.RelatedReference.ReferenceIssuer,
			},
			ReqHdlg: admi_007_001_01.RequestHandling2{
				StsCd:   report.RequestHandling.StatusCode,
				StsDtTm: report.RequestHandling.StatusDateTime,
				Desc:    report.RequestHandling.Description,
			},
		})
	}

	admiDoc := &admi_007_001_01.Document{
		XMLName: xml.Name{Space: "urn:iso:std:iso:20022:tech:xsd:admi.007.001.01", Local: "Document"},
		RctAck: admi_007_001_01.ReceiptAcknowledgementV01{
			MsgId: msgHeader,
			Rpt:   reports,
		},
	}

	return admiDoc, nil
}

func BuildAdmi007(payload []byte, cfg *config.Config) (*admi_007_001_01.Document, error) {
	var message FedNowMessageRctAck
	if err := json.Unmarshal(payload, &message); err != nil {
		return nil, err
	}
	return BuildAdmi007Struct(message, cfg)
}

func ParseAdmi007Struct(admiDoc *admi_007_001_01.Document, appHdr head.BusinessApplicationHeaderV02) (FedNowMessageRctAck, error) {
	rctAck := admiDoc.RctAck

	creationDateTime := common.ISODateTime(appHdr.CreDt)
	if rctAck.MsgId.CreDtTm != nil {
		creationDateTime = common.ISODateTime(*rctAck.MsgId.CreDtTm)
	}

	reports := make([]ReceiptAcknowledgementReport, 0, len(rctAck.Rpt))
	for _, report := range rctAck.Rpt {
		reports = append(reports, ReceiptAcknowledgementReport{
			RelatedReference: ReceiptAcknowledgementReference{
				Reference:       report.RltdRef.Ref,
				MessageName:     report.RltdRef.MsgNm,
				ReferenceIssuer: report.RltdRef.RefIssr,
			},
			RequestHandling: ReceiptAcknowledgementHandling{
				StatusCode:     report.ReqHdlg.StsCd,
				StatusDateTime: report.ReqHdlg.StsDtTm,
				Description:    report.ReqHdlg.Desc,
			},
		})
	}

	fedMsg := FedNowMessageRctAck{
		FedNowMsg: FedNowReceiptAcknowledgement{
			CreationDateTime: creationDateTime,
			Identifier: FedNowIdentifier{
				BusinessMessageID: appHdr.BizMsgIdr,
				MessageType:       appHdr.MsgDefIdr,
				MessageID:         head.Max35Text(rctAck.MsgId.MsgId),
			},
			QueryName: rctAck.MsgId.QryNm,
			Reports:   reports,
		},
	}

	return fedMsg, nil
}
