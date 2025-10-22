package keymanagement

type FedNowKeyExchangeMessage struct {
	PublicKeys []FednowPublicKeys `json:"publicKeys"`
}

type FednowPublicKeys struct {
	KeyStatus               string `json:"keyStatus"`
	StatusDateTime          string `json:"statusDateTime"`
	FedNowStatusDescription string `json:"fedNowStatusDescription"`
	FedNowKeyID             string `json:"fedNowKeyID"`
	Name                    string `json:"name"`
	EncodedPublicKey        string `json:"encodedPublicKey"`
	Encoding                string `json:"encoding"`
	Algorithm               string `json:"algorithm"`
	KeyExpirationDateTime   string `json:"keyExpirationDateTime"`
}
