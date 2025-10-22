package keymanagement

import (
	"github.com/mbanq/iso20022-go/ISO20022/FedNowKeyExchange"
)

func Parse(keyDoc *FedNowKeyExchange.FedNowPublicKeyResponses) (FedNowKeyExchangeMessage, error) {

	var publicKeys []FednowPublicKeys

	for _, pk := range keyDoc.PublicKeys {
		newKey := FednowPublicKeys{
			KeyStatus:               pk.FedNowMessageSignatureKeyStatus.KeyStatus,
			StatusDateTime:          pk.FedNowMessageSignatureKeyStatus.StatusDateTime,
			FedNowStatusDescription: string(pk.FedNowMessageSignatureKeyStatus.FedNowStatusDescription),
			FedNowKeyID:             string(pk.FedNowMessageSignatureKey.FedNowKeyID),
			Name:                    string(pk.FedNowMessageSignatureKey.Name),
			EncodedPublicKey:        pk.FedNowMessageSignatureKey.EncodedPublicKey,
			Encoding:                string(pk.FedNowMessageSignatureKey.Encoding),
			KeyExpirationDateTime:   pk.FedNowMessageSignatureKey.KeyExpirationDateTime,
		}

		if pk.FedNowMessageSignatureKey.Algorithm != nil {
			newKey.Algorithm = string(*pk.FedNowMessageSignatureKey.Algorithm)
		}

		publicKeys = append(publicKeys, newKey)
	}

	keys := FedNowKeyExchangeMessage{
		PublicKeys: publicKeys,
	}

	return keys, nil
}
