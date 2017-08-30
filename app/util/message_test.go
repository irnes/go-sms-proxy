package util

import (
	"strings"
	"testing"
)

func TestCompositeMessaage(t *testing.T) {
	const (
		Numerals     = "0123456789"                                           // 10
		Alphabet     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // 52
		Punctuaction = "~!@#$%^&*()-_+={}[]\\|<,>.?/\"';:`"                   // 32
		ASCII        = Alphabet + Numerals + Punctuaction                     // 94
	)

	testCases := []struct {
		origMessage          string
		expectedParts        int
		expectedFirstPartLen int
		expectedLastPartLen  int
	}{
		{
			origMessage:          ASCII,
			expectedParts:        1,
			expectedFirstPartLen: 94,
			expectedLastPartLen:  0,
		},
		{
			origMessage:          strings.Repeat("@", SingleSMSMaxLen),
			expectedParts:        1,
			expectedFirstPartLen: 160,
			expectedLastPartLen:  0,
		},
		{
			origMessage:          strings.Repeat("@", SingleSMSMaxLen+1),
			expectedParts:        2,
			expectedFirstPartLen: CompositeSMSMaxLen,
			expectedLastPartLen:  8,
		},
		{
			origMessage:          strings.Repeat("@", 2*CompositeSMSMaxLen),
			expectedParts:        2,
			expectedFirstPartLen: CompositeSMSMaxLen,
			expectedLastPartLen:  CompositeSMSMaxLen,
		},
		{
			origMessage:          strings.Repeat("@", 2*CompositeSMSMaxLen+1),
			expectedParts:        3,
			expectedFirstPartLen: CompositeSMSMaxLen,
			expectedLastPartLen:  1,
		},
	}

	for _, tc := range testCases {
		cm := NewCompositeMessage(tc.origMessage)

		if len(cm.MessageParts) != tc.expectedParts {
			t.Fatalf("Expected parts to be %d but got %d",
				tc.expectedParts, len(cm.MessageParts))
		}

		firstPart := cm.MessageParts[0].Body
		if len(firstPart) != tc.expectedFirstPartLen {
			t.Fatalf("Expected length to be %d but got %d",
				tc.expectedFirstPartLen, len(firstPart))
		}

		if tc.expectedLastPartLen > 1 {
			lastIdx := tc.expectedParts - 1
			lastPart := cm.MessageParts[lastIdx:][0].Body
			if len(lastPart) != tc.expectedLastPartLen {
				t.Fatalf("Expected length of last part to be %d but got %d",
					tc.expectedLastPartLen, len(lastPart))
			}
		}
	}

}
