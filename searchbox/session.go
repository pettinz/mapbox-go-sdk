package searchbox

import (
	"crypto/rand"
	"fmt"
)

// NewSessionToken generates a new UUIDv4 session token.
// Session tokens are required for the Suggest/Retrieve workflow
// and should be reused for the entire autocomplete session.
func NewSessionToken() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		// Fallback to a deterministic but unique token if crypto/rand fails
		// This should never happen in practice
		panic("failed to generate session token: " + err.Error())
	}

	// Set version (4) and variant (10) bits for UUIDv4
	b[6] = (b[6] & 0x0f) | 0x40 // Version 4
	b[8] = (b[8] & 0x3f) | 0x80 // Variant 10

	return fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
