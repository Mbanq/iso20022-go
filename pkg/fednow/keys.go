package fednow

import (
	"bytes"
	"encoding/xml"

	"github.com/mbanq/iso20022-go/ISO20022/FedNowKeyExchange"
	"github.com/mbanq/iso20022-go/pkg/fednow/keymanagement"
)

func GetPublicKeys(xmlData []byte) (keymanagement.FedNowKeyExchangeMessage, error) {
	decoder := xml.NewDecoder(bytes.NewReader(xmlData))

	for {
		token, err := decoder.Token()
		if err != nil {
			return keymanagement.FedNowKeyExchangeMessage{}, err
		}

		if se, ok := token.(xml.StartElement); ok {
			if se.Name.Local == "FedNowPublicKeyResponses" {
				var keyResponses FedNowKeyExchange.FedNowPublicKeyResponses
				if err := decoder.DecodeElement(&keyResponses, &se); err != nil {
					return keymanagement.FedNowKeyExchangeMessage{}, err
				}
				return keymanagement.Parse(&keyResponses)
			}
		}
	}
}
