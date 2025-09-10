package fednow

// FedNowMessage is an interface for all FedNow message types.
type FedNowMessage interface {
	IsFedNowMessage()
}
