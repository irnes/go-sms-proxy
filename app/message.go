package app

import ()

// BaseMessage defines an insitial message payload
type BaseMessage struct {
	Recipient  int    `json:"recipient"`
	Originator string `json:"originator"` // max 11 chars
	Message    string `json:"message"`
}
