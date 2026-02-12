package fednow

// FedNowMessage is an interface for all FedNow message types.
type FedNowMessage interface {
	IsFedNowMessage()
}

// WrapperPreferrer is an optional interface that a FedNowMessage can implement
// to indicate which FedNow envelope wrapper element should be used during XML
// generation. This is needed when a single ISO message type (e.g. camt.029.001.09)
// maps to multiple wrapper elements in the FedNow XSD (e.g. FedNowReturnRequestResponse
// vs FedNowInformationRequestResponse). Messages that do not implement this
// interface will fall back to the first matching wrapper found in the XSD.
type WrapperPreferrer interface {
	PreferredWrapper() string
}
