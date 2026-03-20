package layers

import (
	"encoding/base64"
	"fmt"
)

func PresentationEncode(data []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(data))
}

func PresentationDecode(data []byte) []byte {
	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		fmt.Println("Layer 6 (Presentation): decode error:", err)
		return data
	}
	fmt.Println("Layer 6 (Presentation): Base64 decoded")
	return decoded
}
