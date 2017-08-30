package util

import (
	"bytes"
	"encoding/hex"
	"math/rand"
)

const (
	SingleSMSMaxLen    = 160
	CompositeSMSMaxLen = 153
	OriginatorMaxLen   = 11
)

// NewUDH create new User Data Header using provided parameters
// ref is a reference number (must be the same for all parts of the same larger messages)
// total is the total number of SMS messages in the large message
// seq is the sequence number of the current SMS message
func NewUDH(ref, total, seq uint8) UDH {
	udh := UDH{
		0x00,
		0x03,
		byte(ref),
		byte(total),
		byte(seq),
	}

	return udh
}

// UDH represents User Data Header
type UDH [5]byte

// String returns the hexadecimal encoding of UDH prepended with UDH length
func (u UDH) String() string {
	return hex.EncodeToString([]byte{byte(len(u)), u[0], u[1], u[2], u[3], u[4]})
}

// NewMessage create a new instance of Message
func NewMessage(ref, total, seq uint8, body string) *Message {
	udh := NewUDH(ref, total, seq)
	m := &Message{
		Header: udh,
		Body:   body,
	}

	return m
}

// Message struct holds preprocessed SMS payload
type Message struct {
	Header UDH
	Body   string
}

// SplitMessage splits message m to multiple parts whose max length is n
func SplitMessage(m string, n int) []string {
	part := ""
	parts := []string{}

	runes := bytes.Runes([]byte(m))
	l := len(runes)
	for i, r := range runes {
		part = part + string(r)
		if (i+1)%n == 0 {
			parts = append(parts, part)
			part = ""
		} else if (i + 1) == l {
			parts = append(parts, part)
		}
	}

	return parts
}

// NewCompositeMessage creates a new instance of CompositeMessaage
// using provided message payload. It splits a large message into
// multiple parts up to 153 chars leaving 7 chars (49bits) for UDH
// Each message part is prepended fixed an appropriate UDH when sending
func NewCompositeMessage(message string) *CompositeMessaage {
	// generate reference number needed by the receiving end
	// to figure out which part belong to which larhe message
	ref := uint8(rand.Intn(255))

	cm := &CompositeMessaage{
		Ref: ref,
	}

	// split message to multiple parts
	if len(message) > SingleSMSMaxLen {
		mparts := SplitMessage(message, CompositeSMSMaxLen)
		total := uint8(len(mparts))
		for i, mpart := range mparts {
			seq := uint8(i + 1)
			cm.MessageParts = append(cm.MessageParts, NewMessage(cm.Ref, total, seq, mpart))
		}
	} else {
		// there is no need to split message to multiple parts if shorter than 160 chars
		total := uint8(1)
		seq := uint8(1)
		cm.MessageParts = append(cm.MessageParts, NewMessage(cm.Ref, total, seq, message))
	}

	return cm
}

// CompositeMessaage holds parts of concatenated SMS
type CompositeMessaage struct {
	Ref          uint8
	MessageParts []*Message
}
