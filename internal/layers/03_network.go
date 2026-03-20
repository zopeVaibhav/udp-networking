package layers

import (
	"bytes"
	"fmt"
	"osi/internal/models"
)

func NetworkEncode(data []byte, source models.Machine,
	destination models.Machine) []byte {
	header := []byte("SRC:" + source.IPAddress + "|DST:" + destination.IPAddress + "|")
	return append([]byte(header), data...)
}

func NetworkDecode(data []byte) []byte {
	parts := bytes.SplitN(data, []byte("|"), 3)
	fmt.Println("Layer 3 (Network):", string(parts[0]), "|", string(parts[1]))
	return parts[2]
}
