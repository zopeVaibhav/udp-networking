package layers

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func NewSessionID() string {
	buf := make([]byte, 8)
	rand.Read(buf)
	return hex.EncodeToString(buf)
}

func SessionEncode(data []byte, sessionID string) []byte {
	return append([]byte("SID:"+sessionID+"|"), data...)
}

func SessionDecode(data []byte) []byte {
	parts := bytes.SplitN(data, []byte("|"), 2)
	fmt.Println("Layer 5 (Session):", string(parts[0]))
	return parts[1]
}
