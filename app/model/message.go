package model

// BaseMessage defines an initial message payload
// {"recipient":31612345678, "originator":"MessageBird", "message":"This is a test message."}
type BaseMessage struct {
	Recipient  string `json:"recipient"`  // min 9
	Originator string `json:"originator"` // max 11 chars
	Message    string `json:"message"`
}
